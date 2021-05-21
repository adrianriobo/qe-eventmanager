package ocp

import (
	"github.com/mitchellh/mapstructure"

	"github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/messaging/umb"
	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

const (
	buildComplete string = "VirtualTopic.qe.ci.product-scenario.ascerra.build.complete"
	testComplete  string = "VirtualTopic.qe.ci.product-scenario.ascerra.test.complete"
	// testError     string = "VirtualTopic.qe.ci.product-scenario.ascerra.test.error"
)

type ProductScenarioBuild struct {
}

func New() ProductScenarioBuild {
	return ProductScenarioBuild{}
}

func (p ProductScenarioBuild) GetDestination() string {
	return buildComplete
}
func (p ProductScenarioBuild) Handler(event interface{}) error {
	var data BuildComplete

	if err := mapstructure.Decode(event, &data); err != nil {
		return err
	}
	// Business Logic
	for _, product := range data.Artifact.Products {
		if product.Name == "openshift" {
			_, err := pipelines.RunInteropOCP(product.Id, util.GenerateCorrelation())
			if err != nil {
				logging.Error(err)
			}
			// We will take info from status to send back the results
			var response TestComplete
			mockResponse(&data, &response)

			return umb.Send(testComplete, response)
		}
	}
	return nil
}

func mockResponse(source *BuildComplete, response *TestComplete) {
	response.Artifact = source.Artifact
	response.Run = Run{
		URL: "https://crcqe-jenkins-csb-codeready.cloud.paas.psi.redhat.com/view/qe-bundle-baremetal/job/qe/job/bundle_baremetal_macos14-brno/284",
		Log: "https://crcqe-jenkins-csb-codeready.cloud.paas.psi.redhat.com/view/qe-bundle-baremetal/job/qe/job/bundle_baremetal_macos14-brno/284/console"}
	response.Test = Test{
		Category:  "interoperability",
		Namespace: "interop",
		TestType:  "product-scenario",
		Result:    "passed",
		Runtime:   "1800",
		XunitUrls: []string{
			"https://crcqe-jenkins-csb-codeready.cloud.paas.psi.redhat.com/view/qe-bundle-baremetal/job/qe/job/bundle_baremetal_macos14-brno/284/console",
			"https://crcqe-jenkins-csb-codeready.cloud.paas.psi.redhat.com/view/qe-bundle-baremetal/job/qe/job/bundle_baremetal_macos14-brno/284/console"}}
}
