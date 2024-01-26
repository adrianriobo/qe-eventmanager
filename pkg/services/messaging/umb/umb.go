package umb

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/status"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/impl/amqp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/impl/stomp"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

const (
	Stomp string = "stomp"
	Amqp  string = "amqp"
)

var breakingErrors = []string{
	"remote error: tls: user canceled",
	"connection timed out"}

type umbInformation struct {
	consumerID                                   string
	protocol                                     string
	brokers                                      []string
	certificateFile, privateKeyFile, caCertsFile []byte
}

type umb struct {
	umbInformation umbInformation
	client         api.ClientInterface
	subscriptions  map[string]*subscription
	consumers      *sync.WaitGroup
	handlers       *sync.WaitGroup
	breakingError  chan string
	subscribe      sync.Mutex
	send           sync.Mutex
	active         bool
}

type subscription struct {
	topic        string
	subscription api.SubscriptionInterface
	handlers     []api.MessageHandler
	active       bool
}

var _umb *umb

func InitClient(consumerID, protocol string, brokers []string, certificateFile, privateKeyFile, caCertsFile []byte) (err error) {
	_umb, err = initUMB(umbInformation{
		consumerID:      consumerID,
		protocol:        protocol,
		brokers:         brokers,
		certificateFile: certificateFile,
		privateKeyFile:  privateKeyFile,
		caCertsFile:     caCertsFile,
	})
	return
}

func initUMB(umbInfo umbInformation) (umbClient *umb, err error) {
	umbClient = &umb{
		umbInformation: umbInfo,
		subscriptions:  make(map[string]*subscription),
		consumers:      &sync.WaitGroup{},
		handlers:       &sync.WaitGroup{},
		breakingError:  make(chan string),
		active:         true}
	umbClient.client, err = createClient(
		umbInfo.protocol,
		umbInfo.brokers,
		umbInfo.certificateFile,
		umbInfo.privateKeyFile,
		umbInfo.caCertsFile)
	// In case of recconnect it will re create the client and subscriptions
	go umbClient.handleBreakingError()
	return
}

func Send(destination string, message interface{}) error {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		logging.Errorf("Failed to marshal data")
		return err
	}
	return SendBytes(destination, marshalledMessage)
}

func SendBytes(destination string, message []byte) error {
	_umb.send.Lock()
	defer _umb.send.Unlock()
	logging.Debugf("Sending message %s\n, to %s", string(message), destination)
	return _umb.client.Send(destination, message)
}

func Subscribe(subscriptionID, topic string, handlers []api.MessageHandler) error {
	return _umb.subscribeTopic(subscriptionID, topic, handlers)
}

func GracefullShutdown() {
	_umb.active = false
	for _, subscription := range _umb.subscriptions {
		if err := subscription.subscription.Unsubscribe(); err != nil {
			logging.Error(err)
		}
		logging.Infof("Unsubscribing %s", subscription.topic)
	}
	_umb.consumers.Wait()
	_umb.handlers.Wait()
	_umb.client.Disconnect()
	logging.Infof("Client disconnected from UMB")
}

// In case of error on subscription we will force close
// subscription and client an regenate all of them
func (umb *umb) handleBreakingError() {
	<-umb.breakingError
	logging.Debugf("Service will shutdown gracefully")
	if umb.active {
		// GracefullShutdown
		GracefullShutdown()
	}
	// Send signal to mark listerner as unhealthy
	logging.Debugf("Sending signal for unhealthy service")
	status.SendSignal()
}

func (umb *umb) subscribeTopic(subscriptionID, topic string, handlers []api.MessageHandler) error {
	umb.subscribe.Lock()
	defer umb.subscribe.Unlock()
	logging.Infof("Adding a subscription %s on topic %s", subscriptionID, topic)
	internalSubscription, err := umb.client.Subscribe(umbTopic(subscriptionID, topic), handlers)
	if err != nil {
		return err
	}
	umb.subscriptions[subscriptionID] = &subscription{
		topic:        topic,
		subscription: internalSubscription,
		handlers:     handlers,
		active:       true}
	umb.consumers.Add(1)
	go consume(umb.subscriptions[subscriptionID], umb.breakingError)
	// Close the error channel when no consumers left
	go func() {
		umb.consumers.Wait()
		close(umb.breakingError)
	}()
	return nil
}

func createClient(protocol string, brokers []string,
	certificateFile, privateKeyFile, caCertsFile []byte) (api.ClientInterface, error) {
	switch protocol {
	case Stomp:
		return stomp.Create(certificateFile, privateKeyFile, caCertsFile, brokers)
	case Amqp:
		return amqp.Create(certificateFile, privateKeyFile, caCertsFile, brokers)
	default:
		return nil, fmt.Errorf("%s is not supported", protocol)
	}
}

func consume(subscription *subscription, breakingError chan string) {
	defer _umb.consumers.Done()
	for subscription.active {
		msg, err := subscription.subscription.Read()
		if err != nil {
			contains, _ := util.SliceItem(
				breakingErrors,
				func(e string) bool { return strings.Contains(err.Error(), e) },
				func(e string) string { return e })
			if contains {
				// Send cause for reconnect
				breakingError <- fmt.Sprintf("%v on topic %s", err.Error(), subscription.topic)
			}
			logging.Errorf("Error reading from topic: %s. %s", subscription.topic, err)
			break
		}
		logging.Debugf("New message from %s, adding new handler for it", subscription.topic)
		for _, handler := range subscription.handlers {
			_umb.handlers.Add(1)
			go handle(msg, handler)
		}
	}
	logging.Debugf("Finalize consumer for subscription %s", subscription.topic)
}

func handle(msg []byte, handler api.MessageHandler) {
	defer _umb.handlers.Done()
	if err := handler.Match(msg); err == nil {
		if err := handler.Handle(msg); err != nil {
			logging.Errorf("error handling the msg %v", err)
		}
	}
}

// Umb uses identified consumer queues acting as (virtual) topics
// for subscriptions, the full queue name is based on the pattern:
// "Consumer.$SERVICE_ACCOUNT_NAME.$SUBSCRIPTION_ID.VirtualTopic.>"
func umbTopic(subscriptionID, topic string) string {
	topicCrumbs := strings.Split(topic, ".")
	return fmt.Sprintf("Consumer.%s.%s-%s.%s",
		_umb.umbInformation.consumerID,
		subscriptionID,
		topicCrumbs[len(topicCrumbs)-1],
		topic)
}
