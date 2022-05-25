package amqp

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/Azure/go-amqp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/tls"
)

const (
	schema                 = "amqp+ssl"
	defaultConnIdleTimeout = 60 * time.Minute
)

type Client struct {
	Client  *amqp.Client
	Session *amqp.Session
}

type Subscription struct {
	Receiver *amqp.Receiver
}

func Create(certificateFile, privateKeyFile, caCertsFile []byte,
	brokers []string) (api.ClientInterface, error) {
	tlsConfig, err :=
		tls.CreateTLSConfig(certificateFile, privateKeyFile, caCertsFile)
	if err != nil {
		return nil, err
	}
	var client *amqp.Client
	for _, broker := range brokers {
		address := fmt.Sprintf("%s://%s", schema, broker)
		logging.Infof("Connecting to broker %s", address)
		client, err = amqp.Dial(address, amqp.ConnTLSConfig(tlsConfig))
		amqp.ConnIdleTimeout(defaultConnIdleTimeout)
		if err == nil {
			logging.Debugf("Established TCP connection to broker %s", address)
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
	err = s.Receiver.AcceptMessage(context.TODO(), msg)
	if err != nil {
		return nil, err
	}
	if msg.GetData() == nil {
		return []byte(fmt.Sprintf("%v", msg.Value)), nil
	}
	return msg.GetData(), nil
}

func (s Subscription) Unsubscribe() (err error) {
	err = s.Receiver.Close(context.TODO())
	return
}

func (c *Client) Subscribe(destination string,
	handlers []api.MessageHandler) (api.SubscriptionInterface, error) {
	receiver, err := c.Session.NewReceiver(
		amqp.LinkSourceAddress(destination),
		amqp.LinkCredit(10),
	)
	if err != nil {
		return nil, err
	}
	return &Subscription{Receiver: receiver}, nil
}

func (c *Client) Send(destination string, message []byte) error {
	topic := fmt.Sprintf("topic://%s", destination)
	sender, err := c.Session.NewSender(
		amqp.LinkTargetAddress(topic),
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(),
		30*time.Second)
	err = sender.Send(ctx, amqp.NewMessage(message))
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
