package pipelines

import (
	pipelineClientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
)

func NewClient() (*pipelineClientset.Clientset, error) {
	// Initially only allow handle same cluster, the solution should run inside a cluster
	// with SA allowing to create the required resources to handle events
	restconfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := pipelineClientset.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}
	// pipelineRuns := clientset.TektonV1beta1().PipelineRuns("test")
	// pipelineRuns.Create()
	return clientset, nil
}
