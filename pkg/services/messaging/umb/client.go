package umb

import (
	"encoding/json"
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	umb "github.com/adrianriobo/qe-eventmanager/pkg/util/umb"
	stomp "github.com/go-stomp/stomp/v3"
)

const (
	consumerId string = "Consumer.psi-crcqe-openstack.1231231232."
)

var (
	defaultACKMode stomp.AckMode = stomp.AckAuto
)

type Client struct {
	connection    *umb.UMBConnection
	subscriptions []*stomp.Subscription
	consumers     *sync.WaitGroup
	handlers      *sync.WaitGroup
	subscribe     sync.Mutex
	send          sync.Mutex
}

var client Client

func NewClient(certificateFile, privateKeyFile, caCertsFile string, brokers []string) error {
	// Configure
	client.connection = umb.NewUMBConnection(
		certificateFile,
		privateKeyFile,
		caCertsFile,
		brokers)
	// Connect to UMB
	if err := client.connection.Connect(); err != nil {
		return err
	}
	// Initialize waitgroup
	client.consumers = &sync.WaitGroup{}
	client.handlers = &sync.WaitGroup{}
	return nil
}

// TODO add selector based on regex??
func Subscribe(virtualTopic string, handler func(event interface{}) error) error {
	destination := consumerId + virtualTopic
	client.subscribe.Lock()
	defer client.subscribe.Unlock()
	subscription, err := client.connection.FailoverSubscribe(destination, defaultACKMode)
	if err != nil {
		return err
	}
	client.subscriptions = append(client.subscriptions, subscription)
	client.consumers.Add(1)
	go consume(subscription, handler)
	return nil
}

func Send(destination string, message interface{}) error {
	client.send.Lock()
	defer client.send.Unlock()
	return client.connection.FailoverSend("/topic/"+destination, message)
}

func consume(subscription *stomp.Subscription, handler func(event interface{}) error) {
	defer client.consumers.Done()
	for subscription.Active() {
		msg, err := subscription.Read()
		if err != nil {
			if !subscription.Active() {
				break
			}
			logging.Errorf("Error reading from topic: %s. %s", subscription.Destination(), err)
			break
		}
		client.handlers.Add(1)
		go handle(msg, handler)
	}
}

func handle(msg *stomp.Message, handler func(event interface{}) error) {
	// when finish remove from group
	defer client.handlers.Done()
	// heavy consuming may regex over string
	var event interface{}
	// logging.Debugf("Print message %+v", string(msg.Body[:]))
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logging.Error(err)
	}
	if err := handler(event); err != nil {
		logging.Error(err)
	}
}

func GracefullShutdown() {
	for _, subscription := range client.subscriptions {
		if err := subscription.Unsubscribe(); err != nil {
			logging.Error(err)
			// Force consume as finished ?
		}
		client.consumers.Done()
	}
	client.handlers.Wait()
	client.connection.Disconnect()
}
