package rhel

import (
	interopPipelineRHEL "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines/interop-rhel"
	interopEvent "github.com/adrianriobo/qe-eventmanager/pkg/event/interop"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/mitchellh/mapstructure"
)

const (
	topicBuildComplete string = "VirtualTopic.qe.ci.product-scenario.vipatel.build.complete"
	topicTestComplete  string = "VirtualTopic.qe.ci.product-scenario.vipatel.test.complete"

	// topicBuildComplete string = "VirtualTopic.qe.ci.product-scenario.build.complete"
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

func (p ProductScenarioBuild) GetDestination() string {
	return topicBuildComplete
}
func (p ProductScenarioBuild) Handler(event interface{}) error {
	var data BuildComplete

	if err := mapstructure.Decode(event, &data); err != nil {
		return err
	}
	// Business Logic
	var rhelVersion, baseosURL, appstreamURL string
	var codereadyContainersMessage bool = false
	for _, product := range data.Artifact.Products {
		if product.Name == "rhel" {
			rhelVersion = product.Id
			baseosURL, appstreamURL = getRepositoriesURLs(product.Repos)
		}
		if product.Name == "codeready_containers" {
			codereadyContainersMessage = true
		}
	}
	// Filtering this will be improved in future versions
	if len(rhelVersion) > 0 && codereadyContainersMessage {
		name, xunitURL, err :=
			interopPipelineRHEL.Run(rhelVersion, baseosURL, appstreamURL)
		if err != nil {
			logging.Error(err)
		}
		// We will take info from status to send back the results
		response := buildResponse(name, xunitURL, &data)
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

func buildResponse(name, xunitURL string, source *BuildComplete) *TestComplete {
	return &TestComplete{
		Artifact: source.Artifact,
		Run: interopEvent.Run{
			URL: interopPipelineRHEL.GetPipelinerunDashboardUrl(name),
			Log: interopPipelineRHEL.GetPipelinerunDashboardUrl(name)},
		Test: interopEvent.Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    "passed",
			Runtime:   "1800",
			XunitUrls: []string{xunitURL}},
	}
}
