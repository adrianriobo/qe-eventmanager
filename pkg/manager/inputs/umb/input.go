package umb

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/actions"
	"github.com/adrianriobo/qe-eventmanager/pkg/manager/flows"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb/api"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/json"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

// An umb flow will start with a message from umb
// it will create a subscription and a handler
func Add(input flows.UmbInput, action actions.Runnable) error {
	if err := umb.Subscribe(
		input.Topic,
		[]api.MessageHandler{new(input.Filters, action)}); err != nil {
		return err
	}
	return nil
}

type umbFlow struct {
	action  actions.Runnable
	filters []flows.UmbInputFilter
}

func new(filters []flows.UmbInputFilter, action actions.Runnable) umbFlow {
	return umbFlow{
		action:  action,
		filters: filters}
}

func (u umbFlow) Handle(event []byte) error {
	return u.action.Run()
}

func (u umbFlow) Match(event []byte) error {
	var filters []string
	for _, filter := range u.filters {
		filters = append(filters, filter.JsonPath)
	}
	match, err := json.MatchFilters(event, filters)
	if err != nil {
		logging.Errorf("Error checking filters for event %v", err)
	}
	if !match {
		return fmt.Errorf("filters do not match, message will not be processed")
	}
	return nil
}
