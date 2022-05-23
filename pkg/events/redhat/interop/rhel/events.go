package rhel

import (
	"encoding/json"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop"
)

func CreateTestComplete(dahsboardURL, xunitURL,
	duration, resultStatus string, artifactFromEvent []byte, systemFromEvent []byte) (*TestComplete, error) {
	var artifact Artifact
	if err := getNode(artifactFromEvent, artifact); err != nil {
		return nil, err
	}
	var system []interop.System
	if err := getNode(systemFromEvent, artifact); err != nil {
		return nil, err
	}
	return &TestComplete{
		Artifact: artifact,
		Run: interop.Run{
			URL: dahsboardURL,
			Log: dahsboardURL},
		GenerateAt: time.Now().Format(time.RFC3339Nano),
		System:     system,
		Test: interop.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    resultStatus,
			Runtime:   duration,
			XunitUrls: []string{xunitURL}}}, nil
}

func CreateTestError(dahsboardURL string, artifactFromEvent, systemFromEvent []byte) (*TestError, error) {
	var artifact Artifact
	if err := getNode(artifactFromEvent, artifact); err != nil {
		return nil, err
	}
	var system []interop.System
	if err := getNode(systemFromEvent, artifact); err != nil {
		return nil, err
	}
	return &TestError{
		Artifact: artifact,
		Run: interop.Run{
			URL: dahsboardURL,
			Log: dahsboardURL},
		GenerateAt: time.Now().Format(time.RFC3339Nano),
		System:     system,
		Error: interop.Error{
			Reason: "Testing failed due to infrastructure issues",
		},
		Test: interop.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario"}}, nil
}

// func getArtifact(source []byte) (target Artifact, err error) {
// 	err = json.Unmarshal(source, &target)
// 	return
// }

func getNode(source []byte, target interface{}) error {
	return json.Unmarshal(source, &target)
}
