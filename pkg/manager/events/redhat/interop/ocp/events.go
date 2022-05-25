package ocp

import (
	"encoding/json"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/manager/events/redhat/interop"
)

func CreateTestComplete(dahsboardURL, xunitURL,
	duration, resultStatus, contactName, contactEmail string,
	artifactFromEvent []byte, systemFromEvent []byte) (*TestComplete, error) {
	// artifact, err := getArtifact(artifactFromEvent)
	// if err != nil {
	// 	return nil, err
	// }
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
		Contact: interop.Contact{
			Name:  contactName,
			Email: contactEmail,
		},
		Test: interop.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    resultStatus,
			Runtime:   duration,
			XunitUrls: []string{xunitURL}}}, nil
}

func CreateTestError(dahsboardURL, contactName, contactEmail string,
	artifactFromEvent, systemFromEvent []byte) (*TestError, error) {
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
		Contact: interop.Contact{
			Name:  contactName,
			Email: contactEmail,
		},
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
