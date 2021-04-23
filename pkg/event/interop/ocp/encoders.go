package ocp

import (
	"encoding/json"

	"github.com/adrianriobo/qe-eventmanager/pkg/logging"
)

func Unmarshal(source []byte) (*[]Event, error) {
	logging.Debug("Unmarshalling ocp event")
	var event []Event
	if err := json.Unmarshal(source, &event); err != nil {
		return nil, err
	}
	logging.Debugf("Unmarshalled: %+v", event)
	return &event, nil
}

func Marshal(source Event) ([]byte, error) {
	logging.Debug("Marshalling ocp event")
	return json.Marshal(source)
}
