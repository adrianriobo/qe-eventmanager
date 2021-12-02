package umb

import (
	"encoding/json"
	"sync"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	umb "github.com/adrianriobo/qe-eventmanager/pkg/util/umb"
	stomp "github.com/go-stomp/stomp/v3"
)

var (
	defaultACKMode stomp.AckMode = stomp.AckAuto
)

type subscriptionManager struct {
	subscription *stomp.Subscription
	handlers     []func(event interface{}) error
}

type Client struct {
	connection           *umb.UMBConnection
	subscriptionManagers []subscriptionManager
	consumers            *sync.WaitGroup
	handlers             *sync.WaitGroup
	subscribe            sync.Mutex
	send                 sync.Mutex
	active               bool
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
	client.active = true
	return nil
}

func Subscribe(consumerId, virtualTopic string, handlers []func(event interface{}) error) error {
	destination := consumerId + "." + virtualTopic
	return addSubscription(destination, handlers)
}

func Send(destination string, message interface{}) error {
	client.send.Lock()
	defer client.send.Unlock()
	return client.connection.FailoverSend("/topic/"+destination, message)
}

func GracefullShutdown() {
	client.active = false
	for _, subscriptionManager := range client.subscriptionManagers {
		if err := subscriptionManager.subscription.Unsubscribe(); err != nil {
			logging.Error(err)
		}
		logging.Infof("Unsubscribing %s", subscriptionManager.subscription.Destination())
	}
	client.consumers.Wait()
	client.handlers.Wait()
	client.connection.Disconnect()
	logging.Infof("Client disconnected from UMB")
}

func consume(client *Client, subscription *stomp.Subscription, handlers []func(event interface{}) error) {
	defer client.consumers.Done()
	for subscription.Active() {
		msg, err := subscription.Read()
		if err != nil {
			if !subscription.Active() {
				if client.active {
					logging.Debugf("Reconnecting from failing subscription %s", subscription.Destination())
					if err = reconnect(client, subscription.Id()); err != nil {
						logging.Errorf("Error reconnecting from topic: %s. %s", subscription.Destination(), err)
						break
					}
				} else {
					// Should manage if the case is for the gracefull unsubscription or some error
					//...as so a new subscription should be managed
					logging.Debugf("Read message from inactive subscription %s", subscription.Destination())
					break
				}
			}
			logging.Errorf("Error reading from topic: %s. %s", subscription.Destination(), err)
			break
		}
		logging.Infof("New message from %s, adding new handler for it", subscription.Destination())
		for _, handler := range handlers {
			client.handlers.Add(1)
			go handle(client, msg, handler)
		}
	}
	logging.Debugf("Finalize consumer for subscription %s", subscription.Destination())
}

func reconnect(client *Client, subscriptionId string) error {
	subscriptionManager := findSubscription(subscriptionId)
	return addSubscription(subscriptionManager.subscription.Destination(), subscriptionManager.handlers)
}

func addSubscription(destination string, handlers []func(event interface{}) error) error {
	client.subscribe.Lock()
	var subscriptionManager subscriptionManager
	defer client.subscribe.Unlock()
	logging.Infof("Adding a subscription to %s", destination)
	subscription, err := client.connection.FailoverSubscribe(destination, defaultACKMode)
	if err != nil {
		return err
	}
	subscriptionManager.subscription = subscription
	subscriptionManager.handlers = handlers
	client.subscriptionManagers = append(client.subscriptionManagers, subscriptionManager)
	client.consumers.Add(1)
	go consume(&client, subscription, handlers)
	return nil
}

func findSubscription(subscriptionId string) subscriptionManager {
	for _, subscriptionManager := range client.subscriptionManagers {
		if subscriptionManager.subscription.Id() == subscriptionId {
			return subscriptionManager
		}
	}
	return subscriptionManager{}
}

func handle(client *Client, msg *stomp.Message, handler func(event interface{}) error) {
	defer client.handlers.Done()
	// heavy consuming may regex over string, jsonpath
	var event map[string]interface{}
	logging.Debugf("Print message %+v", string(msg.Body[:]))
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logging.Error(err)
	}
	if err := handler(event); err != nil {
		logging.Error(err)
	}
}
