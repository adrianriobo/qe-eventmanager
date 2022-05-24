package forward

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

type ForwardAction struct {
	forwardType string
}

func Create(actionInfo flows.ForwardAction) (*ForwardAction, error) {
	action := &ForwardAction{forwardType: actionInfo.Type}
	return action, nil
}

func (f ForwardAction) Run(event []byte) error {
	logging.Debugf("Got message: %+v", string(event[:]))
	return nil
}
