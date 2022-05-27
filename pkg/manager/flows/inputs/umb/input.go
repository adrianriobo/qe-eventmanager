package umb

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows/actions"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func Add(flowName string, input flows.UMBInput, action actions.Runnable) error {
	if err := umb.Subscribe(
		input.Topic,
		[]api.MessageHandler{new(flowName, input.Filters, action)}); err != nil {
		return err
	}
	return nil
}

type umbFlow struct {
	flowName string
	action   actions.Runnable
	filters  []string
}

func new(flowName string, filters []string, action actions.Runnable) umbFlow {
	return umbFlow{
		flowName: flowName,
		action:   action,
		filters:  filters}
}

func (u umbFlow) Handle(event []byte) error {
	return u.action.Run(event)
}

func (u umbFlow) Match(event []byte) error {
	var filters []string
	filters = append(filters, u.filters...)
	match, err := json.MatchFilters(event, filters)
	if err != nil {
		logging.Errorf("Error checking filters for event %v", err)
	}
	if !match {
		return fmt.Errorf("filters do not match, message will not be processed")
	}
	logging.Debugf("Found event marching the filters for flow %s", u.flowName)
	return nil
}
