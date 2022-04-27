package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func CreateTLSConfig(cert_path, key_path, ca_path string) (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(cert_path, key_path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load ssl cert/key from: %s, %s",
			cert_path, key_path)
	}
	logging.Debugf("Cert and key loaded successfully")

	// Load CA cert
	caCert, err := ioutil.ReadFile(ca_path)
	if err != nil {
		return nil, fmt.Errorf("failed to load ssl CA bundle: %s", ca_path)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		// ca certs is now well set
		// InsecureSkipVerify: true,
	}, nil
}
