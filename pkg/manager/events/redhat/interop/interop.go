package interop

import (
	"encoding/json"
)

func AdaptEventNodes(artifactFromEvent, systemFromEvent []byte,
	artifact any) ([]System, error) {
	var system []System
	if err := json.Unmarshal(systemFromEvent, &system); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(artifactFromEvent, artifact); err != nil {
		return nil, err
	}
	return system, nil
}
