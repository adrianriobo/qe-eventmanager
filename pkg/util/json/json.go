package json

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/spyzhov/ajson"
)

func MatchFilters(event []byte, filters []string) (bool, error) {
	root, err := ajson.Unmarshal(event)
	if err != nil {
		return false, err
	}
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

func GetStringValue(event []byte, jsonPath string) (string, error) {
	root, err := ajson.Unmarshal(event)
	if err != nil {
		return "", err
	}
	nodes, err := root.JSONPath(jsonPath)
	if err != nil {
		return "", fmt.Errorf("error with %v", err)
	}
	if len(nodes) != 1 {
		return "", fmt.Errorf("error with %v", err)
	}
	if value, err := nodes[0].GetString(); err != nil {
		return "", fmt.Errorf("error with %v", err)
	} else {
		return value, nil
	}
}
