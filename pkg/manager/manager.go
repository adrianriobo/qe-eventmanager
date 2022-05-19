package manager

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/actions"
	actionForward "github.com/adrianriobo/qe-eventmanager/pkg/manager/actions/forward"
	actionTekton "github.com/adrianriobo/qe-eventmanager/pkg/manager/actions/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	inputsUMB "github.com/adrianriobo/qe-eventmanager/pkg/manager/inputs/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/providers"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/cicd/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/file"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func Initialize(providersFilePath string, flowsFilePath []string) {
	providers, flows, err := loadFiles(providersFilePath, flowsFilePath)
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
	if err := manageFlows(flows); err != nil {
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

func loadFiles(providersFilePath string, flowsFilePath []string) (*providers.Providers, *[]flows.Flow, error) {
	var structuredProviders providers.Providers
	var structuredFlows []flows.Flow
	if len(providersFilePath) > 0 {
		if err := file.LoadFileAsStruct(providersFilePath, &structuredProviders); err != nil {
			logging.Errorf("Can not load providers file: %v", err)
			return nil, nil, err
		}
	}
	if len(flowsFilePath) > 0 {
		for _, flowFilePath := range flowsFilePath {
			if len(flowFilePath) > 0 {
				var flow flows.Flow
				if err := file.LoadFileAsStruct(flowFilePath, &flow); err != nil {
					// Should try to keep loading remaining rules if exist
					logging.Errorf("Can not load flows file: %v", err)
					return nil, nil, err
				}
				structuredFlows = append(structuredFlows, flow)
			}
		}
	}
	return &structuredProviders, &structuredFlows, nil
}

func createTektonClient(info providers.Tekton) (err error) {
	var workspaces []tektonClient.WorkspaceBinding
	if len(info.Workspaces) > 0 {
		for _, item := range info.Workspaces {
			var adaptedItem tektonClient.WorkspaceBinding
			adaptedItem.Name = item.Name
			adaptedItem.PVC = item.PVC
			workspaces = append(workspaces, adaptedItem)
		}
	}
	kubeconfig := []byte("")
	if len(info.Kubeconfig) > 0 {
		kubeconfig, err = base64.StdEncoding.DecodeString(info.Kubeconfig)
		if err != nil {
			return
		}
	}
	return tektonClient.CreateClient(kubeconfig, info.Namespace, workspaces, info.ConsoleURL)
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

// For each flow defined:
// Create action that action
func manageFlows(flows *[]flows.Flow) error {
	if len(*flows) > 0 {
		for _, flow := range *flows {
			logging.Debugf("Setting up flow: %v", flow)
			action, err := getAction(flow)
			if err != nil {
				logging.Errorf("Find error with flow %s:%v", flow.Name, err)
				break
			}
			err = addActionToInput(flow, action)
			if err != nil {
				logging.Errorf("Find error with flow %s:%v", flow.Name, err)
				break
			}
		}
	}
	return nil
}

func getAction(flow flows.Flow) (actions.Runnable, error) {
	// if flow.Action.TektonPipelineAction != nil {
	if !util.IsEmpty(flow.Action.TektonPipeline) {
		//Create the action
		action, err := actionTekton.Create(flow.Action.TektonPipeline)
		if err != nil {
			return nil, err
		}
		return action, nil
	}
	if !util.IsEmpty(flow.Action.Forward) {
		//Create the action
		action, err := actionForward.Create(flow.Action.Forward)
		if err != nil {
			return nil, err
		}
		return action, nil
	}
	return nil, fmt.Errorf("action is invalid")
}

func addActionToInput(flow flows.Flow, action actions.Runnable) error {
	// if flow.Input.UmbInput != nil {
	if !util.IsEmpty(flow.Input.UMB) {
		inputsUMB.Add(flow.Input.UMB, action)
	}
	return nil
}
