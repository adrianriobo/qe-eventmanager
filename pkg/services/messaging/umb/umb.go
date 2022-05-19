package umb

import (
	"fmt"
	"strings"
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/impl/amqp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/impl/stomp"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

const (
	Stomp string = "stomp"
	Amqp  string = "amqp"
)

type umb struct {
	consumerID    string
	client        api.ClientInterface
	subscriptions []subscription
	consumers     *sync.WaitGroup
	handlers      *sync.WaitGroup
	subscribe     sync.Mutex
	send          sync.Mutex
	active        bool
}

type subscription struct {
	topic        string
	subscription api.SubscriptionInterface
	handlers     []api.MessageHandler
	active       bool
}

var _umb *umb

func CreateClient(consumerID, protocol string, brokers []string, certificateFile, privateKeyFile, caCertsFile []byte) (err error) {
	_umb = &umb{
		consumerID: consumerID,
		consumers:  &sync.WaitGroup{},
		handlers:   &sync.WaitGroup{},
		active:     true}
	_umb.client, err = createClient(protocol, brokers, certificateFile, privateKeyFile, caCertsFile)
	return
}

func Send(destination string, message interface{}) error {
	_umb.send.Lock()
	defer _umb.send.Unlock()
	return _umb.client.Send(destination, message)
}

func Subscribe(topic string, handlers []api.MessageHandler) error {
	_umb.subscribe.Lock()
	defer _umb.subscribe.Unlock()
	logging.Infof("Adding a subscription to %s", topic)
	internalSubscription, err := _umb.client.Subscribe(umbTopic(topic), handlers)
	if err != nil {
		return err
	}
	var subscription = subscription{
		topic:        topic,
		subscription: internalSubscription,
		handlers:     handlers,
		active:       true}
	_umb.subscriptions = append(_umb.subscriptions, subscription)
	_umb.consumers.Add(1)
	go consume(&subscription)
	return nil
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

func consume(subscription *subscription) {
	defer _umb.consumers.Done()
	for subscription.active {
		msg, err := subscription.subscription.Read()
		if err != nil {
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
			logging.Error(err)
		}
	}
}

// Umb uses identified consumer queues acting as (virtual) topics
// for subscriptions, the full queue name is based on the pattern:
// "Consumer.$SERVICE_ACCOUNT_NAME.$SUBSCRIPTION_ID.VirtualTopic.>"
func umbTopic(topic string) string {
	subscriptionId := strings.Split(topic, ".")
	return fmt.Sprintf("Consumer.%s.%s.%s",
		_umb.consumerID,
		subscriptionId[len(subscriptionId)-1],
		topic)
}
