package actions

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/configuration/flows"
	actionForward "github.com/adrianriobo/qe-eventmanager/pkg/manager/flows/actions/forward"
	actionTekton "github.com/adrianriobo/qe-eventmanager/pkg/manager/flows/actions/tekton"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
)

type Runnable interface {
	Run(event []byte) error
}

func CreateAction(flow flows.Flow) (Runnable, error) {
	// if flow.Action.PipelineAction != nil {
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
