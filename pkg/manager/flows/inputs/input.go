package inputs

import (
	"github.com/devtools-qe-incubator/eventmanager/pkg/configuration/flows"
	"github.com/devtools-qe-incubator/eventmanager/pkg/manager/flows/actions"
	inputACK "github.com/devtools-qe-incubator/eventmanager/pkg/manager/flows/inputs/ack"
	inputsUMB "github.com/devtools-qe-incubator/eventmanager/pkg/manager/flows/inputs/umb"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util"
)

func AddActionToInput(flow flows.Flow, action actions.Runnable) error {
	var ack inputACK.ACK
	if !util.IsEmpty(flow.Input.ACK) {
		//return ack as action to run from input
		ack = inputACK.CreateACK(flow.Input.ACK)
	}
	if !util.IsEmpty(flow.Input.UMB) {
		return inputsUMB.Add(flow.Name, flow.Input.UMB, ack, action)
	}
	return nil
}
