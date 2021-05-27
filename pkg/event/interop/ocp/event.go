package ocp

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

const (
	topicBuildComplete string = "VirtualTopic.qe.ci.product-scenario.ascerra.build.complete"
	topicTestComplete  string = "VirtualTopic.qe.ci.product-scenario.ascerra.test.complete"
	// testError     string = "VirtualTopic.qe.ci.product-scenario.ascerra.test.error"
)

var (
	serversids []string = []string{"macos14-brno", "macos15-brno", "windows10-brno", "rhel8-brno"}
	platforms  []string = []string{"fedora33", "rhel79", "rhel83"}
	files      []string = []string{"basic.xml", "config.xml", "story_health.xml",
		"story_marketplace.xml", "story_registry.xml", "cert_rotation.xml",
		"proxy.xml", "integration.xml"}
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
	for _, product := range data.Artifact.Products {
		if product.Name == "openshift" {
			name, correlation, _, err :=
				pipelines.RunInteropOCP(product.Id, util.GenerateCorrelation(),
					strings.Join(serversids[:], ","),
					strings.Join(platforms[:], ","))
			if err != nil {
				logging.Error(err)
			}
			// We will take info from status to send back the results
			response := buildResponse(name, correlation, &data)
			return umb.Send(topicTestComplete, response)
		}
	}
	return nil
}

func buildResponse(name, correlation string, source *BuildComplete) *TestComplete {
	return &TestComplete{
		Artifact: source.Artifact,
		Run: Run{
			URL: pipelines.GetPipelinerunDashboardUrl(name),
			Log: pipelines.GetPipelinerunDashboardUrl(name)},
		Test: Test{
			Category:  "interoperability",
			Namespace: "interop",
			TestType:  "product-scenario",
			Result:    "passed",
			Runtime:   "1800",
			XunitUrls: xunitFilesUrls(correlation)},
	}
}

func xunitFilesUrls(correlation string) []string {
	var xunitUrls []string
	datalakeUrl := "http://10.0.110.220:9000/logs"
	t := time.Now().Local()
	logsDate := fmt.Sprint(t.Format("20060102"))
	servers := append(serversids, platforms...)
	for _, server := range servers {
		for _, file := range files {
			url := fmt.Sprintf("%s/%s/%s/%s/%s",
				datalakeUrl, logsDate, correlation, server, file)
			xunitUrls = append(xunitUrls, url)
		}
	}
	return xunitUrls
}
