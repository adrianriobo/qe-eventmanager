package rhel

import "github.com/adrianriobo/qe-eventmanager/pkg/events/redhat/interop"

func CreateTestComplete(eventSchema, dahsboardURL, xunitURL,
	duration, resultStatus string, artifact Artifact) TestComplete {
	return TestComplete{
		Artifact: artifact,
		Run: interop.Run{
			URL: dahsboardURL,
			Log: dahsboardURL},
		Test: interop.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    resultStatus,
			Runtime:   duration,
			XunitUrls: []string{xunitURL}}}
}
