package tekton

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/adrianriobo/qe-eventmanager/pkg/events"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/cicd/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	tektonUtil "github.com/adrianriobo/qe-eventmanager/pkg/util/tekton"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonAction struct {
	actionInfo flows.TektonPipelineAction
}

func Create(actionInfo flows.TektonPipelineAction) (*TektonAction, error) {
	action := &TektonAction{actionInfo: actionInfo}
	return action, nil
}

func (a TektonAction) Run(event []byte) error {
	// Check if pipelinerun require some parameter from event
	pipelineRunParameters, err := parsePipelineParameters(
		a.actionInfo.Params, event)
	if err != nil {
		return err
	}
	// Create the pipelinerun spec with params from event
	pipelineRunSpec := createPipelineRun(
		a.actionInfo.Name, pipelineRunParameters)
	logging.Debugf("Creating pipelinerun spec: %v", *pipelineRunSpec)
	// Use tekton client to create the run
	pipelineRun, err := tektonClient.ApplyPipelinerun(pipelineRunSpec)
	if err != nil {
		return err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go tektonClient.AddInformer(pipelineRun.GetName(), status, informerStopper)
	return manageResults(<-status, pipelineRun.GetName(),
		event, a.actionInfo.Success, a.actionInfo.Error)
}

func createPipelineRun(pipelineName string, params []v1beta1.Param) *v1beta1.PipelineRun {
	pipelineRunName := fmt.Sprintf("%s-", pipelineName)
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineName},
			Params:      params},
	}
}

// Pending json from event
func parsePipelineParameters(pipelineFlowParams []flows.TektonPipelineParam,
	event []byte) (parameters []v1beta1.Param, err error) {
	for _, flowParam := range pipelineFlowParams {
		var value string
		if len(flowParam.JsonPath) > 0 {
			value, err = json.GetStringValue(event, flowParam.JsonPath)
			if err != nil {
				return
			}
		} else if len(flowParam.Value) > 0 {
			value = flowParam.Value
		} else {
			return nil, fmt.Errorf("parameter %s require a value, review flow definition", flowParam.Name)
		}
		param := v1beta1.Param{
			Name:  flowParam.Name,
			Value: *v1beta1.NewArrayOrString(value)}
		parameters = append(parameters, param)
	}
	return
}

func manageResults(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success, errorFlow flows.Error) error {
	if tektonUtil.IsSuccessful(status) {
		return manageSuccess(status, pipelineRunName, event, success)
	}
	return manageError(pipelineRunName, event, errorFlow)
}

func manageSuccess(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	event []byte, success flows.Success) error {
	if success.UMB.EventSchema == events.RedHatInteropOCPTestComplete ||
		success.UMB.EventSchema == events.RedHatInteropRHELTestComplete {
		//Default jsonpath could be passed as fields
		artifactNode, err := json.GetNodeAsByteArray([]byte(event), "$.artifact")
		if err != nil {
			return err
		}
		systemNode, err := json.GetNodeAsByteArray([]byte(event), "$.system")
		if err != nil {
			return err
		}
		dashboardURL := tektonClient.GetPipelinerunDashboardUrl(pipelineRunName)
		xunitURLs, duration, resultStatus :=
			getPipelineRunResults(status, pipelineRunName, success)
		response, err := events.GenerateRedHatInteropTestComplete(success.UMB.EventSchema,
			dashboardURL, xunitURLs, duration, resultStatus, artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(success.UMB.Topic, response)
	}
	return nil
}

func manageError(pipelineRunName string,
	event []byte, errorFLow flows.Error) error {
	if errorFLow.UMB.EventSchema == events.RedHatInteropOCPTestError ||
		errorFLow.UMB.EventSchema == events.RedHatInteropRHELTestError {
		//Default jsonpath could be passed as fields
		artifactNode, err := json.GetNodeAsByteArray([]byte(event), "$.artifact")
		if err != nil {
			return err
		}
		systemNode, err := json.GetNodeAsByteArray([]byte(event), "$.system")
		if err != nil {
			return err
		}
		dashboardURL := tektonClient.GetPipelinerunDashboardUrl(pipelineRunName)
		response, err := events.GenerateRedHatInteropTestError(errorFLow.UMB.EventSchema,
			dashboardURL, artifactNode, systemNode)
		if err != nil {
			return err
		}
		return umb.Send(errorFLow.UMB.Topic, response)
	}
	return nil
}

func getPipelineRunResults(status *v1beta1.PipelineRunStatus, pipelineRunName string,
	success flows.Success) (xunitURLs, duration, resultStatus string) {
	xunitURLs = tektonUtil.GetResultValue(status.PipelineResults,
		getPipelineResultItem(success, events.RedHatInteropXunitURL))
	duration = tektonUtil.GetResultValue(status.PipelineResults,
		getPipelineResultItem(success, events.RedHatInteropDuration))
	resultStatus = tektonUtil.GetResultValue(status.PipelineResults,
		getPipelineResultItem(success, events.RedHatInteropResultStatus))
	return
}

func getPipelineResultItem(success flows.Success, item string) string {
	idx := slices.IndexFunc(success.UMB.EventFields,
		func(e flows.UMBEventField) bool { return e.Name == item })
	if idx == -1 {
		return ""
	}
	return success.UMB.EventFields[idx].PipelineResultName
}
