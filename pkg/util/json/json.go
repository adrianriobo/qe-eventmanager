package json

import (
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/spyzhov/ajson"
)

func matchFilters(event string, filters []string) (bool, error) {
	json := []byte(event)
	root, _ := ajson.Unmarshal(json)
	for _, filter := range filters {
		node, err := root.JSONPath(filter)
		if err != nil {
			logging.Error("Error with %v", err)
			return false, err
		}
		if len(node) == 0 {
			logging.Error("Error event does not match the filters")
			return false, nil
		}
		logging.Debug("Found event marching the filters")
	}
	return true, nil
}
