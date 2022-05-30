package tekton

import (
	"fmt"
	"strings"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows/results"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/cicd/tekton"
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

// Action run implementation
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
	// Use tekton client to create the run
	pipelineRun, err := tektonClient.ApplyPipelinerun(pipelineRunSpec)
	if err != nil {
		return err
	}
	logging.Debugf("Created pipelinerun : %v", pipelineRun)
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	logging.Debugf("Added informer for pipelinerun %s", pipelineRun.GetName())
	go tektonClient.AddInformer(pipelineRun.GetName(), status, informerStopper)
	return results.ManageResults(<-status, pipelineRun.GetName(),
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

func parsePipelineParameters(pipelineFlowParams []flows.NameValuePair,
	event []byte) (parameters []v1beta1.Param, err error) {
	for _, flowParam := range pipelineFlowParams {
		var value string
		value, err = getParamValue(flowParam.Value, event)
		if err != nil {
			break
		}
		param := v1beta1.Param{
			Name:  flowParam.Name,
			Value: *v1beta1.NewArrayOrString(value)}
		parameters = append(parameters, param)
	}
	return
}

func getParamValue(valueExpression string, event []byte) (string, error) {
	// Check if value should be picked from pipelinerun results
	if strings.HasPrefix(valueExpression, json.JSONPathPreffix) {
		return json.GetStringValue(event, valueExpression)
	}
	// Otherwise value is a constant value
	return valueExpression, nil
}
