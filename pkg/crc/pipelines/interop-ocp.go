package pipelines

import (
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

func Run(ocpVersion, correlation string) (*v1beta1.PipelineRunStatus, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcNamespace, getSpec(ocpVersion, correlation))
	if err != nil {
		return nil, err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	defer close(status)
	if err := pipelines.AddInformer(crcNamespace, pipelinerun.GetName(), status); err != nil {
		return nil, err
	}
	result := <-status
	return result, nil
}

func getSpec(ocpVersion, correlation string) *v1beta1.PipelineRun {
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
