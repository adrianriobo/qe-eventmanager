package stomp

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/message/stomp"
	stompv3 "github.com/go-stomp/stomp/v3"
)

var (
	defaultACKMode stompv3.AckMode = stompv3.AckAuto
)

type Client struct {
	Connection *stomp.Connection
}

type Subscription struct {
	Subscription *stompv3.Subscription
}

func Create(certificateFile, privateKeyFile, caCertsFile []byte,
	brokers []string) (api.ClientInterface, error) {
	connection, err := stomp.NewConnection(
		certificateFile,
		privateKeyFile,
		caCertsFile,
		brokers)
	if err != nil {
		return nil, err
	}
	var stompClient = Client{
		Connection: connection}
	if err := stompClient.Connection.Connect(); err != nil {
		return nil, err
	}
	return &stompClient, nil
}

func (s Subscription) Read() ([]byte, error) {
	msg, err := s.Subscription.Read()
	if err != nil {
		return nil, err
	}
	return msg.Body, nil
}

func (s Subscription) Unsubscribe() (err error) {
	err = s.Subscription.Unsubscribe()
	return
}

func (c *Client) Send(destination string, message interface{}) error {
	return c.Connection.FailoverSend("/topic/"+destination, message)
}

func (c *Client) Disconnect() {
	c.Connection.Disconnect()
	logging.Infof("Client disconnected from UMB")
}

func (c *Client) Subscribe(destination string,
	handlers []api.MessageHandler) (api.SubscriptionInterface, error) {
	subscription, err := c.Connection.FailoverSubscribe(destination, defaultACKMode)
	if err != nil {
		return nil, err
	}
	return Subscription{Subscription: subscription}, nil
}
