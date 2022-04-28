package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/Azure/go-amqp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/tls"
)

type Client struct {
	Client  *amqp.Client
	Session *amqp.Session
}

// https://stackoverflow.com/questions/67491806/how-do-you-connect-to-an-amqp-1-0-topic-not-queue-in-golang
type Subscription struct {
	Receiver *amqp.Receiver
}

func Create(certificateFile, privateKeyFile, caCertsFile string, brokers []string) (api.ClientInterface, error) {
	tlsConfig, err := tls.CreateTLSConfig(certificateFile, privateKeyFile, caCertsFile)
	if err != nil {
		return nil, err
	}
	var client *amqp.Client
	for _, url := range brokers {
		logging.Debugf("Connecting to broker %s", url)
		client, err = amqp.Dial("", amqp.ConnTLSConfig(tlsConfig))
		if err == nil {
			logging.Debugf("Established TCP connection to broker %s", url)
			break
		}
		// log.WithField("broker", url).Warning("Connection to broker failed: %s", err.Error())
		logging.Debugf("Connection to broker failed: %v", err)
	}
	if client == nil {
		return nil, fmt.Errorf("unable to establish connection for provided brokers")
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:  client,
		Session: session}, err
}

func (s Subscription) Read() ([]byte, error) {
	msg, err := s.Receiver.Receive(context.TODO())
	if err != nil {
		return nil, err
	}
	return msg.GetData(), nil
}

func (s Subscription) Unsubscribe() (err error) {
	err = s.Receiver.Close(context.TODO())
	return
}

func (c *Client) Subscribe(destination string, handlers []func(event interface{}) error) (api.SubscriptionInterface, error) {
	receiver, err := c.Session.NewReceiver(
		amqp.LinkSourceAddress(destination),
		amqp.LinkCredit(10),
	)
	if err != nil {
		return nil, err
	}
	return &Subscription{Receiver: receiver}, nil
}

func (c *Client) Send(destination string, message interface{}) error {
	umbTopic := fmt.Sprintf("/topic/%s", destination)
	sender, err := c.Session.NewSender(
		amqp.LinkTargetAddress(umbTopic),
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(),
		5*time.Second)
	jsonData, err := json.Marshal(message)
	if err != nil {
		logging.Errorf("Failed to marshal data")
		cancel()
		return err
	}
	err = sender.Send(ctx, amqp.NewMessage(jsonData))
	if err != nil {
		cancel()
		return err
	}
	if err := sender.Close(ctx); err != nil {
		logging.Error("Error closing the amqp sender")
	}
	cancel()
	return nil
}

func (c *Client) Disconnect() {
	if err := c.Client.Close(); err != nil {
		logging.Error("Error closing amqp client connetion")
	}
}
