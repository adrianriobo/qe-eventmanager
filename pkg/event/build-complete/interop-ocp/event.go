package interopOCP

import (
	interopPipelineOCP "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines/interop-ocp"
	buildComplete "github.com/adrianriobo/qe-eventmanager/pkg/event/build-complete"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/mitchellh/mapstructure"
)

const (
	topicTestComplete string = "VirtualTopic.qe.ci.product-scenario.test.complete"
	// testError     string = "VirtualTopic.qe.ci.product-scenario.test.error"
)

type ProductScenarioBuild struct {
}

func New() ProductScenarioBuild {
	return ProductScenarioBuild{}
}

func (p ProductScenarioBuild) Handler(event interface{}) error {
	var data BuildComplete

	if err := mapstructure.Decode(event, &data); err != nil {
		return err
	}
	// Business Logic
	var openshiftVersion string = ""
	var codereadyContainersMessage bool = false
	for _, product := range data.Artifact.Products {
		if product.Name == "openshift" {
			openshiftVersion = product.Id
		}
		if product.Name == "codeready_containers" {
			codereadyContainersMessage = true
		}
	}
	// Filtering this will be improved in future versions
	if len(openshiftVersion) > 0 && codereadyContainersMessage {
		name, xunitURL, duration, resultStatus, err :=
			interopPipelineOCP.Run(openshiftVersion)
		if err != nil {
			logging.Error(err)
		}
		// We will take info from status to send back the results
		response := buildResponse(name, xunitURL, duration, resultStatus, &data)
		return umb.Send(topicTestComplete, response)
	}
	return nil
}

func buildResponse(name, xunitURL, duration, resultStatus string, source *BuildComplete) *TestComplete {
	return &TestComplete{
		Artifact: source.Artifact,
		Run: buildComplete.Run{
			URL: interopPipelineOCP.GetPipelinerunDashboardUrl(name),
			Log: interopPipelineOCP.GetPipelinerunDashboardUrl(name)},
		Test: buildComplete.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    resultStatus,
			Runtime:   duration,
			XunitUrls: []string{xunitURL}},
	}
}
