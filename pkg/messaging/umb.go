/*
Package UMB provides a simplified way to connect to Unified Message Bus (UMB)
using SSL client certs and failver over STOMP protocol.

Primary reasons for its existence:
 * Enforced 'application/json' content type
 * SSL cert loading and TLS connection is a bit convoluted
 * go-stomp does not support failover for JBoss AMQ servers

Ideally it should probably be mostly merged into go-stomp one day

*/

package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/logging"
	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type UMBConnection struct {
	// count of retry until Failover* function fail
	FailoverRetryCount uint
	// delay between reconnects
	FailoverRetryDelay float32
	// real delay = retry delay counter * FailoverRetryDelay
	// first retry is without delay and next delays are longer and longer

	ssl_cert_path string
	ssl_key_path  string
	ssl_ca_path   string
	hosts         []string
	tlsConfig     *tls.Config
	tlsConn       *tls.Conn
	stompConn     *stomp.Conn
	connOpts      []func(*stomp.Conn) error
}

const (
	DefaultFailoverRetryCount = 1
	DefaultFailoverRetryDelay = 0.1
)

func NewUMBConnection(sslCertPath, sslKeyPath, sslCaPath string, hosts []string, connOpts ...func(*stomp.Conn) error) *UMBConnection {
	return &UMBConnection{
		ssl_cert_path:      sslCertPath,
		ssl_key_path:       sslKeyPath,
		ssl_ca_path:        sslCaPath,
		hosts:              hosts,
		connOpts:           connOpts,
		FailoverRetryCount: DefaultFailoverRetryCount,
		FailoverRetryDelay: DefaultFailoverRetryDelay,
	}
}

func (c *UMBConnection) FailoverSend(destination string, body interface{}, opts ...func(*frame.Frame) error) error {
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

func (c *UMBConnection) FailoverSubscribe(destination string, ack stomp.AckMode, opts ...func(*frame.Frame) error) (*stomp.Subscription, error) {
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

func (c *UMBConnection) Connect() error {
	if c.tlsConfig == nil {
		err := c.createTLSConfig()
		if err != nil {
			return errors.New("Failed to load SSL certificates: " + err.Error())
		}

	}

	var conn *tls.Conn
	var err error

	for _, url := range c.hosts {
		logging.Debugf("Connecting to broker")
		conn, err = tls.Dial("tcp", url, c.tlsConfig)

		if err == nil {
			logging.Infof("Established TCP connection to broker")
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

func (c *UMBConnection) Disconnect() {
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

func (c *UMBConnection) createTLSConfig() error {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(c.ssl_cert_path, c.ssl_key_path)
	if err != nil {
		return errors.New("Failed to load ssl cert/key from: " + c.ssl_cert_path + ", " + c.ssl_key_path)
	}
	logging.Debugf("Cert and key loaded successfully")

	// Load CA cert
	caCert, err := ioutil.ReadFile(c.ssl_ca_path)
	if err != nil {
		return errors.New("Failed to load ssl CA bundle: " + c.ssl_ca_path)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	c.tlsConfig = tlsConfig
	return nil
}
