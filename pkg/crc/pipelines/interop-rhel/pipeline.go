package interopRHEL

import (
	"fmt"

	crcPipelines "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/http"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/xunit"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRefName string = "interop-rhel"
	pipelineRunName string = pipelineRefName + "-"

	rhelVersionParamName   string = "rhel-version"
	repoBaseosParamName    string = "repo-baseos-url"
	repoAppStreamParamName string = "repo-appstream-url"
	imageIDParamName       string = "image-id"

	xunitURLResultName   string = "results-url"
	qeDurationResultName string = "qe-duration"

	resultStatusPassed string = "passed"
	resultStatusFailed string = "failed"
)

func Run(rhelVersion, repoBaseos, repoAppStream, imageID string) (string, string, string, string, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcPipelines.Namespace, getSpec(rhelVersion, repoBaseos, repoAppStream, imageID))
	if err != nil {
		return "", "", "", "", err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go pipelines.AddInformer(crcPipelines.Namespace, pipelinerun.GetName(), status, informerStopper)
	runStatus := <-status
	xunitURL := getResultValue(runStatus.PipelineResults, xunitURLResultName)
	return pipelinerun.GetName(),
		xunitURL,
		getResultValue(runStatus.PipelineResults, qeDurationResultName),
		getResultState(xunitURL),
		nil
}

// TODO make general vailable
func getResultValue(results []v1beta1.PipelineRunResult, resultParamID string) string {
	for _, result := range results {
		if result.Name == resultParamID {
			return result.Value
		}
	}
	return ""
}

// TODO this should be moved to result parameter from the pipeline
func getResultState(url string) string {
	file, err := http.GetFile(url)
	if err != nil {
		logging.Error(err)
		return ""
	}
	count, err := xunit.CountFailures(file)
	if err != nil {
		logging.Error(err)
		return ""
	}
	if count == 0 {
		return resultStatusPassed
	}
	return resultStatusFailed
}

func getSpec(rhelVersion, repoBaseos, repoAppStream, imageID string) *v1beta1.PipelineRun {
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName, Namespace: crcPipelines.Namespace},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineRefName},
			Params: []v1beta1.Param{
				{Name: rhelVersionParamName, Value: *v1beta1.NewArrayOrString(rhelVersion)},
				{Name: repoBaseosParamName, Value: *v1beta1.NewArrayOrString(repoBaseos)},
				{Name: repoAppStreamParamName, Value: *v1beta1.NewArrayOrString(repoAppStream)},
				{Name: imageIDParamName, Value: *v1beta1.NewArrayOrString(imageID)}},
			Timeout:    &crcPipelines.DefaultTimeout,
			Workspaces: []v1beta1.WorkspaceBinding{crcPipelines.Workspace}},
	}
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf(crcPipelines.DashboardUrlFormat, crcPipelines.DashboardBaseUrl, crcPipelines.Namespace, pipelinerunName)
}
