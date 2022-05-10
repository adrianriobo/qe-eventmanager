package json

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/spyzhov/ajson"
)

func MatchFiltersAsString(event string, filters []string) (bool, error) {
	return MatchFilters([]byte(event), filters)
}

func MatchFilters(event []byte, filters []string) (bool, error) {
	root, _ := ajson.Unmarshal(event)
	for _, filter := range filters {
		node, err := root.JSONPath(filter)
		if err != nil {
			return false, fmt.Errorf("error with %v", err)
		}
		if len(node) == 0 {
			return false, fmt.Errorf("error event does not match the filters")
		}
		logging.Debug("Found event marching the filters")
	}
	return true, nil
}
