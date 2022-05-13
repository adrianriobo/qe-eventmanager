package tekton

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/ci/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
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
		a.actionInfo.PipelineParams, event)
	if err != nil {
		return err
	}
	// Create the pipelinerun spec with params from event
	pipelineRunSpec := createPipelineRun(
		a.actionInfo.PipelineName, pipelineRunParameters)
	logging.Debugf("Creating pipelinerun spec: %v", pipelineRunSpec)
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
	runStatus := <-status
	logging.Debugf("Got the pipelinestatus %v", runStatus)

	// xunitURL := tektonUtil.GetResultValue(runStatus.PipelineResults, xunitURLResultName)
	// return pipelinerun.GetName(),
	// 	xunitURL,
	// 	tektonUtil.GetResultValue(runStatus.PipelineResults, qeDurationResultName),
	// 	tektonUtil.GetResultState(xunitURL),
	// 	nil
	return nil
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
			return nil, fmt.Errorf("Parameter %s require a value, review flow definition", flowParam.Name)
		}
		param := v1beta1.Param{
			Name:  flowParam.Name,
			Value: *v1beta1.NewArrayOrString(value)}
		parameters = append(parameters, param)
	}
	return
}
