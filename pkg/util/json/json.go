package json

import (
	"fmt"

	"github.com/spyzhov/ajson"
)

const (
	JSONPathPreffix = "$."
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
			return false, nil
		}
	}
	return true, nil
}

func GetStringValue(event []byte, jsonPath string) (string, error) {
	node, err := getNode(event, jsonPath)
	if err != nil {
		return "", err
	}
	if node.IsNumeric() {
		value, err := node.GetNumeric()
		if err != nil {
			return "", fmt.Errorf("error with %v", err)
		}
		return fmt.Sprintf("%f", value), nil
	}
	if value, err := node.GetString(); err != nil {
		return "", fmt.Errorf("error with %v", err)
	} else {
		return value, nil
	}
}

func GetNodeAsByteArray(event []byte, jsonPath string) ([]byte, error) {
	node, err := getNode(event, jsonPath)
	if err != nil {
		return nil, err
	}
	return ajson.Marshal(node)
}

func getNode(event []byte, jsonPath string) (*ajson.Node, error) {
	root, err := ajson.Unmarshal(event)
	if err != nil {
		return nil, err
	}
	nodes, err := root.JSONPath(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error with %v", err)
	}
	if len(nodes) != 1 {
		return nil, fmt.Errorf("error with %v", err)
	}
	return nodes[0], nil
}
