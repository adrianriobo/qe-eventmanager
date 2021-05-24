package pipelines

import (
	"fmt"

	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRefName string = "interop-ocp"
	pipelineRunName string = pipelineRefName + "-"

	ocpVersionParamName  string = "ocp-version"
	correlationParamName string = "correlation"
)

func RunInteropOCP(ocpVersion, correlation string) (string, string, *v1beta1.PipelineRunStatus, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcNamespace, getSpecInteropOCP(ocpVersion, correlation))
	if err != nil {
		return "", "", nil, err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go pipelines.AddInformer(crcNamespace, pipelinerun.GetName(), status, informerStopper)
	return pipelinerun.GetName(), correlation, <-status, nil
}

func getSpecInteropOCP(ocpVersion, correlation string) *v1beta1.PipelineRun {
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName, Namespace: crcNamespace},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineRefName},
			Params: []v1beta1.Param{
				{Name: ocpVersionParamName, Value: *v1beta1.NewArrayOrString(ocpVersion)},
				{Name: correlationParamName, Value: *v1beta1.NewArrayOrString(correlation)}},
			Timeout:    &defaultTimeout,
			Workspaces: []v1beta1.WorkspaceBinding{crcWorkspace}},
	}
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf(pipelinesDashboardUrlFormat, pipelinesDashboardBaseUrl, crcNamespace, pipelinerunName)
}
