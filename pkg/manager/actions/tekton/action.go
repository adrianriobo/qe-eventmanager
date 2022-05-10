package tekton

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	tektonClient "github.com/adrianriobo/qe-eventmanager/pkg/services/ci/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonAction struct {
	actionManagerID string
	actionInfo      flows.TektonPipelineAction
}

func Create(actionManagerID string, actionInfo flows.TektonPipelineAction) (*TektonAction, error) {
	var action TektonAction
	action.actionManagerID = actionManagerID
	action.actionInfo = actionInfo
	return &action, nil
}

func (a TektonAction) Run() error {
	pipelineRunSpec := a.createPipelineRun()
	logging.Debugf("Creating pipelinerun spec: %v", pipelineRunSpec)
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

func (a TektonAction) createPipelineRun() *v1beta1.PipelineRun {
	var params []v1beta1.Param
	for _, tuple := range a.actionInfo.PipelineParams {
		param := v1beta1.Param{Name: tuple.Name, Value: *v1beta1.NewArrayOrString(tuple.Value)}
		params = append(params, param)
	}
	pipelineRunName := fmt.Sprintf("%s-%s-", a.actionManagerID, a.actionInfo.PipelineName)
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: a.actionInfo.PipelineName},
			Params:      params},
	}
}
