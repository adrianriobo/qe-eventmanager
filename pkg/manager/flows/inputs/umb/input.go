package umb

import (
	"fmt"

	"github.com/devtools-qe-incubator/eventmanager/pkg/configuration/flows"
	"github.com/devtools-qe-incubator/eventmanager/pkg/manager/flows/actions"
	"github.com/devtools-qe-incubator/eventmanager/pkg/manager/flows/inputs/ack"
	"github.com/devtools-qe-incubator/eventmanager/pkg/services/messaging/umb"
	"github.com/devtools-qe-incubator/eventmanager/pkg/services/messaging/umb/api"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/json"
	"github.com/devtools-qe-incubator/eventmanager/pkg/util/logging"
)

func Add(flowName string, input flows.UMBInput, ack ack.ACK, action actions.Runnable) error {
	if err := umb.Subscribe(
		flowName,
		input.Topic,
		[]api.MessageHandler{new(flowName, input.Filters, ack, action)}); err != nil {
		return err
	}
	return nil
}

type umbFlow struct {
	flowName string
	ack      ack.ACK
	action   actions.Runnable
	filters  []string
}

func new(flowName string, filters []string, ack ack.ACK, action actions.Runnable) umbFlow {
	return umbFlow{
		flowName: flowName,
		ack:      ack,
		action:   action,
		filters:  filters}
}

func (u umbFlow) Handle(event []byte) error {
	if u.ack != nil {
		err := u.ack(event)
		if err != nil {
			logging.Errorf(" error when ack the message %v", err)
		}
	}
	return u.action.Run(event)
}

func (u umbFlow) Match(event []byte) error {
	var filters []string
	filters = append(filters, u.filters...)
	match, err := json.MatchFilters(event, filters)
	if err != nil {
		return err
	}
	if !match {
		logging.Debug("no matching filters, dimissing message")
		return fmt.Errorf("filters do not match, message will not be processed")
	}
	logging.Debugf("Found event marching the filters for flow %s", u.flowName)
	return nil
}
