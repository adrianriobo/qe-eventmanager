package stomp

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	utilTLS "github.com/adrianriobo/qe-eventmanager/pkg/util/tls"
	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type Connection struct {
	FailoverRetryCount uint
	FailoverRetryDelay float32
	hosts              []string
	tlsConfig          *tls.Config
	tlsConn            *tls.Conn
	stompConn          *stomp.Conn
	connOpts           []func(*stomp.Conn) error
}

const (
	DefaultFailoverRetryCount = 1
	DefaultFailoverRetryDelay = 0.1
)

func NewConnection(sslCertPath, sslKeyPath, sslCaPath []byte,
	hosts []string, connOpts ...func(*stomp.Conn) error) (*Connection, error) {
	tlsConfig, err := utilTLS.CreateTLSConfig(sslCertPath, sslKeyPath, sslCaPath)
	if err != nil {
		return nil, err
	}
	return &Connection{
		tlsConfig:          tlsConfig,
		hosts:              hosts,
		connOpts:           connOpts,
		FailoverRetryCount: DefaultFailoverRetryCount,
		FailoverRetryDelay: DefaultFailoverRetryDelay,
	}, nil
}

func (c *Connection) Connect() (err error) {
	var conn *tls.Conn
	for _, url := range c.hosts {
		logging.Infof("Connecting to broker")
		conn, err = tls.Dial("tcp", url, c.tlsConfig)

		if err == nil {
			logging.Infof("Established TCP connection to broker %s", url)
			break
		}
		// log.WithField("broker", url).Warning("Connection to broker failed: %s", err.Error())
	}
	if err != nil {
		return errors.New("Failed tcp broker connection: " + err.Error())
	}

	// disable heartbeat since it actually just makes us disconnect
	// see: https://github.com/go-stomp/stomp/issues/32
	// remove or set to short time to test failover :-)
	opts := append(c.connOpts, stomp.ConnOpt.HeartBeat(0, 0))
	stompConn, err := stomp.Connect(conn, opts...)
	if err != nil {
		return errors.New("Failed stomp connection: " + err.Error())
	}
	c.tlsConn = conn
	c.stompConn = stompConn
	return nil
}

func (c *Connection) FailoverSend(destination string, body interface{}, opts ...func(*frame.Frame) error) error {
	var retryCount uint

	jsonData, err := json.Marshal(body)
	if err != nil {
		logging.Errorf("Failed to marshal data")
		return err
	}

	opts = append(opts, stomp.SendOpt.NoContentLength)
	for retryCount = 0; retryCount <= c.FailoverRetryCount; retryCount++ {
		err = c.stompConn.Send(destination, "application/json", jsonData, opts...)
		if err == nil {
			break
		}

		if retryCount > 0 {
			logging.Debugf("Failed send over stomp right after reconnecting. Permission problems?")
		}
		c.Disconnect()

		time.Sleep(time.Duration(c.FailoverRetryDelay*float32(retryCount)) * time.Second)

		err := c.Connect()
		if err != nil {
			logging.Errorf("Failed to connect to any broker during reconnect: " + err.Error())
		}
	}

	return err
}

func (c *Connection) FailoverSubscribe(destination string, ack stomp.AckMode, opts ...func(*frame.Frame) error) (*stomp.Subscription, error) {
	var (
		sub        *stomp.Subscription
		err        error
		retryCount uint
	)

	for retryCount = 0; retryCount <= c.FailoverRetryCount; retryCount++ {
		sub, err = c.stompConn.Subscribe(destination, ack, opts...)
		if err == nil {
			break
		}

		if retryCount > 0 {
			logging.Debugf("Failed subscribe over stomp right after reconnecting. Permission problems?")
		}
		c.Disconnect()

		time.Sleep(time.Duration(c.FailoverRetryDelay*float32(retryCount)) * time.Second)

		err = c.Connect()
		if err != nil {
			logging.Errorf("Failed to connect to any broker during reconnect: " + err.Error())
		}
	}

	return sub, err
}

func (c *Connection) Disconnect() {
	if c.stompConn != nil {
		err := c.stompConn.Disconnect()
		if err != nil {
			logging.Errorf("error disconnecting")
		}
	}
	if c.tlsConn != nil {
		c.tlsConn.Close()
	}
}
