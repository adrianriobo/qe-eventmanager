package manager

import (
	"os"
	"os/signal"

	"github.com/adrianriobo/qe-eventmanager/pkg/event/interop/ocp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	umb "github.com/adrianriobo/qe-eventmanager/pkg/util/umb"
)

type Client struct {
	certificateFile string
	privateKeyFile  string
	caCertsFile     string
	brokers         []string
}

func New(certificateFile, privateKeyFile, caCertsFile, kubeconfigPath string, brokers []string) *Client {
	// Start pipeline client
	if err := pipelines.NewClient(kubeconfigPath); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	// Start umb client
	return &Client{
		certificateFile: certificateFile,
		privateKeyFile:  privateKeyFile,
		caCertsFile:     caCertsFile,
		brokers:         brokers,
	}
}

func (c Client) Run() {
	connection := umb.NewUMBConnection(
		c.certificateFile,
		c.privateKeyFile,
		c.caCertsFile,
		c.brokers)
	if err := connection.Connect(); err != nil {
		logging.Error(err)
		os.Exit(0)
	}
	productScenarioBuild := ocp.New(connection)
	productScenarioBuild.Init()
	waitForStop()
	// Consumers routine should end gracefully to avoid data losing
	// Handlers routines generated should end gracefully to avoid data losing
	productScenarioBuild.Finish()
	connection.Disconnect()
	logging.Info("Event manager was gracefully stopped. Enjoy your day!")
	os.Exit(0)
}

func waitForStop() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s
}
