package tls

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func CreateTLSConfig(userCertificate, userKey, certificateAuthority []byte) (*tls.Config, error) {
	// Load client cert
	cert, err := tls.X509KeyPair(userCertificate, userKey)
	if err != nil {
		return nil, err
	}
	logging.Debugf("Cert and key loaded successfully")

	// Load CA cert
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(certificateAuthority)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		// ca certs is now well set
		// InsecureSkipVerify: true,
	}, nil
}
