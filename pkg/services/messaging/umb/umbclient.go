package umb

import (
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

type umbSubscription struct {
	subscriptionID string
	handlers       []func(event interface{}) error
}

type umbClient struct {
	umbclient     umbclientInterface
	subscriptions []umbSubscription
	consumers     *sync.WaitGroup
	handlers      *sync.WaitGroup
	subscribe     sync.Mutex
	send          sync.Mutex
	active        bool
}

var clientV2 umbClient

func NewUmbClient(certificateFile, privateKeyFile, caCertsFile string, brokers []string) error {
	clientV2.consumers = &sync.WaitGroup{}
	clientV2.handlers = &sync.WaitGroup{}
	clientV2.active = true
	return nil
}

func (u *umbClient) Send(destination string, message interface{}) error {
	u.send.Lock()
	defer u.send.Unlock()
	return u.umbclient.Send(destination, message)
}

func (u *umbClient) Subscribe(consumerId, virtualTopic string, handlers []func(event interface{}) error) error {
	destination := consumerId + "." + virtualTopic
	return u.addSubscription(destination, handlers)
}

func (u *umbClient) addSubscription(destination string, handlers []func(event interface{}) error) error {
	u.subscribe.Lock()
	var subsciption umbSubscription
	defer u.subscribe.Unlock()
	logging.Infof("Adding a subscription to %s", destination)
	err := u.umbclient.Subscribe(destination, handlers)
	if err != nil {
		return err
	}
	subsciption.subscriptionID = destination
	subsciption.handlers = handlers
	u.subscriptions = append(u.subscriptions, subsciption)
	u.consumers.Add(1)
	return nil
}

type umbclientInterface interface {
	Subscribe(destination string, handlers []func(event interface{}) error) error
	Send(destination string, message interface{}) error
	GracefullShutdown()
}

type umbStompClient struct {
	client *umbClient
	//Move old client + stomp utils
}
type umbAmqpClient struct {
	client *umbClient
	// https://github.com/Azure/go-amqp/blob/main/example_test.go
}
