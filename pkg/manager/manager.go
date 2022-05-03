package manager

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"encoding/base64"

	interopRHEL "github.com/adrianriobo/qe-eventmanager/pkg/event/build-complete/interop-rhel"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/providers"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/rules"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/file"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func Initialize(providersFilePath string, rulesFilePath []string) {
	providers, rules, err := loadFiles(providersFilePath, rulesFilePath)
	if err != nil {
		logging.Errorf("%v", err)
		os.Exit(1)
	}
	if err := createTektonClient(providers.Tekton); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	if err := createUMBClient(providers.UMB); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	if err := manageRules(rules); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	// Execute until stop signal
	waitForStop()
	stop()
	os.Exit(0)
}

func waitForStop() {
	s := make(chan os.Signal, 1)
	signal.Notify(s,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-s
}

func stop() {
	umb.GracefullShutdown()
	logging.Info("Event manager was gracefully stopped. Enjoy your day!")
}

func loadFiles(providersFilePath string, rulesFilePath []string) (*providers.Providers, *[]rules.Rule, error) {
	var structuredProviders providers.Providers
	var structuredRules []rules.Rule
	if len(providersFilePath) > 0 {
		if err := file.LoadFileAsStruct(providersFilePath, &structuredProviders); err != nil {
			logging.Errorf("Can not load providers file: %v", err)
			return nil, nil, err
		}
	}
	if len(rulesFilePath) > 0 {
		for _, ruleFilePath := range rulesFilePath {
			if len(ruleFilePath) > 0 {
				var rule rules.Rule
				if err := file.LoadFileAsStruct(ruleFilePath, &rule); err != nil {
					// Should try to keep loading remaining rules if exist
					logging.Errorf("Can not load rules file: %v", err)
					return nil, nil, err
				}
				structuredRules = append(structuredRules, rule)
			}
		}
	}
	return &structuredProviders, &structuredRules, nil
}

func createTektonClient(info providers.Tekton) (err error) {
	var workspaces []tekton.WorkspaceBinding
	if len(info.Workspaces) > 0 {
		for _, item := range info.Workspaces {
			var adaptedItem tekton.WorkspaceBinding
			adaptedItem.Name = item.Name
			adaptedItem.PVC = item.PVC
			workspaces = append(workspaces, adaptedItem)
		}
	}

	kubeconfig, err := base64.StdEncoding.DecodeString(info.Kubeconfig)
	if err != nil {
		logging.Debugf("%s", string(kubeconfig))
		err = tekton.CreateClient(kubeconfig, info.Namespace, workspaces)
	}
	return
}

func createUMBClient(info providers.UMB) (err error) {
	userCertificate, err :=
		base64.StdEncoding.DecodeString(info.UserCertificate)
	if err != nil {
		return
	}
	userKey, err :=
		base64.StdEncoding.DecodeString(info.UserKey)
	if err != nil {
		return
	}
	certificateAuthority, err :=
		base64.StdEncoding.DecodeString(info.CertificateAuthority)
	if err != nil {
		return
	}
	err = umb.CreateClient(
		info.ConsumerID,
		info.Driver,
		strings.Split(info.Brokers, ","),
		userCertificate,
		userKey,
		certificateAuthority)
	return
}

func manageRules(rules *[]rules.Rule) error {
	if len(*rules) > 0 {
		logging.Debugf("Printing rules content: %v", rules)
		for _, rule := range *rules {
			// Check if input is umb, for the moment the only accepted input
			if err := umb.Subscribe(rule.Input.UmbInput.Topic, []func(event interface{}) error{
				func(event interface{}) error { return interopRHEL.New().Handler(event) }}); err != nil {
				umb.GracefullShutdown()
				return err
			}
		}
	}
	return nil
}
