package interopRHEL

import (
	"fmt"

	crcPipelines "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRefName string = "interop-rhel"
	pipelineRunName string = pipelineRefName + "-"

	rhelVersionParamName   string = "rhel-version"
	repoBaseosParamName    string = "repo-baseos-url"
	repoAppStreamParamName string = "repo-appstream-url"

	xunitURLResultName string = "results-url"
)

func Run(rhelVersion, repoBaseos, repoAppStream string) (string, string, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcPipelines.Namespace, getSpec(rhelVersion, repoBaseos, repoAppStream))
	if err != nil {
		return "", "", err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go pipelines.AddInformer(crcPipelines.Namespace, pipelinerun.GetName(), status, informerStopper)
	runStatus := <-status
	return pipelinerun.GetName(), getResultsURL(runStatus.PipelineResults), nil
}

func getResultsURL(results []v1beta1.PipelineRunResult) string {
	for _, result := range results {
		if result.Name == xunitURLResultName {
			return result.Value
		}
	}
	return ""
}

func getSpec(rhelVersion, repoBaseos, repoAppStream string) *v1beta1.PipelineRun {
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName, Namespace: crcPipelines.Namespace},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineRefName},
			Params: []v1beta1.Param{
				{Name: rhelVersionParamName, Value: *v1beta1.NewArrayOrString(rhelVersion)},
				{Name: repoBaseosParamName, Value: *v1beta1.NewArrayOrString(repoBaseos)},
				{Name: repoAppStreamParamName, Value: *v1beta1.NewArrayOrString(repoAppStream)}},
			Timeout:    &crcPipelines.DefaultTimeout,
			Workspaces: []v1beta1.WorkspaceBinding{crcPipelines.Workspace}},
	}
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf(crcPipelines.DashboardUrlFormat, crcPipelines.DashboardBaseUrl, crcPipelines.Namespace, pipelinerunName)
}
