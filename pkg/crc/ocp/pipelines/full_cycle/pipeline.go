package full_cycle

import (
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	versioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	tektoninformers "github.com/tektoncd/pipeline/pkg/client/informers/externalversions/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cache "k8s.io/client-go/tools/cache"
)

const (
	// TODO move to ENV or file with defaults
	// ocpPreviewMirror string = "http://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview"
	pipelineRefName string = "full-cycle"
	pipelineRunName string = pipelineRefName + "-"
	crcOCPNamespace string = "codeready-container"

	ocpMirrorParamName  string = "ocp-mirror"
	ocpVersionParamName string = "ocp-version"
)

func GetSpec(ocpVersion, ocpMirror *string) *pipelinev1beta1.PipelineRun {
	return &pipelinev1beta1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pipelineRunName,
			Namespace:    crcOCPNamespace,
		},
		Spec: pipelinev1beta1.PipelineRunSpec{
			PipelineRef: &pipelinev1beta1.PipelineRef{
				Name: pipelineRefName,
			},
			Params: []pipelinev1beta1.Param{
				{
					Name:  ocpMirrorParamName,
					Value: *pipelinev1beta1.NewArrayOrString(*ocpMirror),
				},
				{
					Name:  ocpVersionParamName,
					Value: *pipelinev1beta1.NewArrayOrString(*ocpVersion),
				}},
		},
	}
	// TODO create with clientset
}

func SetInformer(clientset versioned.Interface) {
	stopper := make(chan struct{})
	defer close(stopper)
	// https://github.com/kubernetes-client/java/issues/725
	informer := tektoninformers.NewFilteredPipelineRunInformer(clientset, crcOCPNamespace, 0, nil, nil)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{UpdateFunc: func(oldObj, newObj interface{}) {
		pipelineRun := newObj.(*pipelinev1beta1.PipelineRun)
		if pipelineRun.IsDone() {
			// DESIGN best approach one informer per run or one informer and some async mechanism from there
			// when we get the status result on the generated pipelinerun we can close the informer
			close(stopper)
		}
	}})
	informer.Run(stopper)
}
