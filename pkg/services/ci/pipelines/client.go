package pipelines

import (
	"context"
	"fmt"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	informers "github.com/tektoncd/pipeline/pkg/client/informers/externalversions/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var client *clientset.Clientset

// Initially only allow handle same cluster, the solution should run inside a cluster
// with SA allowing to create the required resources to handle events
// var config
func NewClient(kubeconfigPath string) error {
	var config *rest.Config
	var err error
	if kubeconfigPath != "" {
		// use the current context in kubeconfig
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath); err != nil {
			return err
		}
	} else {
		if config, err = rest.InClusterConfig(); err != nil {
			return err
		}
	}
	if client, err = clientset.NewForConfig(config); err != nil {
		return err
	}
	return nil
}

func CreatePipelinerun(namespace string, spec *v1beta1.PipelineRun) (*v1beta1.PipelineRun, error) {
	if err := checkInitialization(); err != nil {
		return nil, err
	}
	return client.TektonV1beta1().PipelineRuns(namespace).Create(context.Background(), spec, v1.CreateOptions{})
}

// DESIGN best approach one informer per run or one informer and some async mechanism from there
// when we get the status result on the generated pipelinerun we can close the informer
func AddInformer(namespace, pipelinerunName string, status chan *v1beta1.PipelineRunStatus) error {
	if err := checkInitialization(); err != nil {
		return err
	}
	informerStopper := make(chan struct{})
	defer close(informerStopper)
	// https://github.com/kubernetes-client/java/issues/725
	informer := informers.NewFilteredPipelineRunInformer(client, namespace, 0, nil, nil)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{UpdateFunc: func(oldObj, newObj interface{}) {
		pipelineRun := newObj.(*v1beta1.PipelineRun)
		if pipelineRun.GetName() == pipelinerunName && pipelineRun.IsDone() {
			// Send the status of the pipelinerun when is done
			status <- &pipelineRun.Status
			close(informerStopper)
		}
	}})
	informer.Run(informerStopper)
	return nil
}

func checkInitialization() error {
	if client == nil {
		return fmt.Errorf("pipelines client is not initialized")
	}
	return nil
}
