package interopOCP

import (
	"fmt"

	crcPipelines "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"
	commonPipelines "github.com/adrianriobo/qe-eventmanager/pkg/util/pipelines"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRefName string = "interop-ocp"
	pipelineRunName string = pipelineRefName + "-"

	ocpVersionParamName string = "ocp-version"

	xunitURLResultName   string = "results-url"
	qeDurationResultName string = "qe-duration"
)

func Run(ocpVersion string) (string, string, string, string, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcPipelines.Namespace, getSpec(ocpVersion))
	if err != nil {
		return "", "", "", "", err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go pipelines.AddInformer(crcPipelines.Namespace, pipelinerun.GetName(), status, informerStopper)
	runStatus := <-status
	xunitURL := commonPipelines.GetResultValue(runStatus.PipelineResults, xunitURLResultName)
	return pipelinerun.GetName(),
		xunitURL,
		commonPipelines.GetResultValue(runStatus.PipelineResults, qeDurationResultName),
		commonPipelines.GetResultState(xunitURL),
		nil
}

func getSpec(ocpVersion string) *v1beta1.PipelineRun {
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName, Namespace: crcPipelines.Namespace},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineRefName},
			Params: []v1beta1.Param{
				{Name: ocpVersionParamName, Value: *v1beta1.NewArrayOrString(ocpVersion)}},
			Timeout:    &crcPipelines.DefaultTimeout,
			Workspaces: []v1beta1.WorkspaceBinding{crcPipelines.Workspace}},
	}
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf(crcPipelines.DashboardUrlFormat, crcPipelines.DashboardBaseUrl, crcPipelines.Namespace, pipelinerunName)
}
