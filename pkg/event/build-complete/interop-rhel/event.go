package interopRHEL

import (
	interopPipelineRHEL "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines/interop-rhel"
	buildComplete "github.com/adrianriobo/qe-eventmanager/pkg/event/build-complete"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/mitchellh/mapstructure"
)

const (
	topicTestComplete string = "VirtualTopic.qe.ci.product-scenario.vipatel.test.complete"

	// topicTestComplete  string = "VirtualTopic.qe.ci.product-scenario.test.complete"
	// testError     string = "VirtualTopic.qe.ci.product-scenario.ascerra.test.error"

	baseosRepositoryName    string = "baseos"
	appstreamRepositoryName string = "appstream"
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
	var rhelVersion, baseosURL, appstreamURL, imageID string
	var codereadyContainersMessage bool = false
	for _, product := range data.Artifact.Products {
		if product.Name == "rhel" {
			rhelVersion = product.Id
			baseosURL, appstreamURL = getRepositoriesURLs(product.Repos)
			logging.Debugf("Got repos baseos: %s, appstream %s", baseosURL, appstreamURL)
			imageID = product.Image
		}
		if product.Name == "codeready_containers" {
			codereadyContainersMessage = true
		}
	}
	// Filtering this will be improved in future versions
	if len(rhelVersion) > 0 && codereadyContainersMessage {
		name, xunitURL, duration, resultStatus, err :=
			interopPipelineRHEL.Run(rhelVersion, baseosURL, appstreamURL, imageID)
		if err != nil {
			logging.Error(err)
		}
		// We will take info from status to send back the results
		response := buildResponse(name, xunitURL, duration, resultStatus, &data)
		return umb.Send(topicTestComplete, response)
	}
	return nil
}

func getRepositoriesURLs(repositories []Repository) (baseosURL, appstreamURL string) {
	for _, repository := range repositories {
		switch repository.Name {
		case baseosRepositoryName:
			baseosURL = repository.BaseUrl
		case appstreamRepositoryName:
			appstreamURL = repository.BaseUrl
		}
	}
	return
}

func buildResponse(name, xunitURL, duration, resultStatus string, source *BuildComplete) *TestComplete {
	return &TestComplete{
		Artifact: source.Artifact,
		Run: buildComplete.Run{
			URL: interopPipelineRHEL.GetPipelinerunDashboardUrl(name),
			Log: interopPipelineRHEL.GetPipelinerunDashboardUrl(name)},
		Test: buildComplete.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    resultStatus,
			Runtime:   duration,
			XunitUrls: []string{xunitURL}},
	}
}
