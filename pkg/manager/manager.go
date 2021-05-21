package manager

import (
	"os"
	"os/signal"

	eventInteropOCP "github.com/adrianriobo/qe-eventmanager/pkg/event/interop/ocp"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func Initialize(certificateFile, privateKeyFile, caCertsFile, kubeconfigPath string, brokers []string) {
	// Start pipeline client
	if err := pipelines.NewClient(kubeconfigPath); err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	// Start umb client
	if err := umb.NewClient(certificateFile, privateKeyFile, caCertsFile, brokers); err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	// Handle events
	if err := handleEvents(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	// Execute until stop signal
	waitForStop()
	stop()
	os.Exit(0)
}

func handleEvents() error {
	eventDefinition := eventInteropOCP.New()
	if err := umb.Subscribe(eventDefinition.GetDestination(), func(event interface{}) error { return eventDefinition.Handler(event) }); err != nil {
		umb.GracefullShutdown()
		return err
	}
	return nil
}

func waitForStop() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s
}

func stop() {
	umb.GracefullShutdown()
	logging.Info("Event manager was gracefully stopped. Enjoy your day!")
}