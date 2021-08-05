package interop

import (
	"fmt"

	crcPipelines "github.com/adrianriobo/qe-eventmanager/pkg/crc/pipelines"
	"github.com/adrianriobo/qe-eventmanager/pkg/services/ci/pipelines"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pipelineRefName string = "interop-ocp"
	pipelineRunName string = pipelineRefName + "-"

	ocpVersionParamName  string = "ocp-version"
	correlationParamName string = "correlation"
	serversidsParamName  string = "servers-ids"
	platformsParamName   string = "platforms"
)

func RunInteropOCP(ocpVersion, correlation, serversids, platforms string) (string, string, *v1beta1.PipelineRunStatus, error) {
	pipelinerun, err := pipelines.CreatePipelinerun(crcPipelines.Namespace, getSpecInteropOCP(ocpVersion, correlation, serversids, platforms))
	if err != nil {
		return "", "", nil, err
	}
	status := make(chan *v1beta1.PipelineRunStatus)
	informerStopper := make(chan struct{})
	defer close(status)
	defer close(informerStopper)
	go pipelines.AddInformer(crcPipelines.Namespace, pipelinerun.GetName(), status, informerStopper)
	return pipelinerun.GetName(), correlation, <-status, nil
}

func getSpecInteropOCP(ocpVersion, correlation, serversids, platforms string) *v1beta1.PipelineRun {
	return &v1beta1.PipelineRun{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{GenerateName: pipelineRunName, Namespace: crcPipelines.Namespace},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{Name: pipelineRefName},
			Params: []v1beta1.Param{
				{Name: ocpVersionParamName, Value: *v1beta1.NewArrayOrString(ocpVersion)},
				{Name: correlationParamName, Value: *v1beta1.NewArrayOrString(correlation)},
				{Name: serversidsParamName, Value: *v1beta1.NewArrayOrString(serversids)},
				{Name: platformsParamName, Value: *v1beta1.NewArrayOrString(platforms)}},
			Timeout:    &crcPipelines.DefaultTimeout,
			Workspaces: []v1beta1.WorkspaceBinding{crcPipelines.Workspace}},
	}
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf(crcPipelines.DashboardUrlFormat, crcPipelines.DashboardBaseUrl, crcPipelines.Namespace, pipelinerunName)
}
