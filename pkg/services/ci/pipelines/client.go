package pipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
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
func AddInformer(namespace, pipelinerunName string, status chan *v1beta1.PipelineRunStatus, informerStopper chan struct{}) {
	if err := checkInitialization(); err != nil {
		logging.Error(err)
	}
	// https://github.com/kubernetes-client/java/issues/725
	informer := informers.NewPipelineRunInformer(client, namespace, 2*time.Minute, nil)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{UpdateFunc: func(oldObj, newObj interface{}) {
		pipelineRun, ok := newObj.(*v1beta1.PipelineRun)
		if !ok {
			logging.Error("error formatting pipelinerun")
		}
		logging.Debugf("Change on pipelinerun %s", pipelineRun.GetName())
		if pipelineRun.GetName() == pipelinerunName &&
			(pipelineRun.IsDone() || pipelineRun.IsCancelled() || pipelineRun.IsTimedOut()) {
			// Send the status of the pipelinerun when is done
			status <- &pipelineRun.Status
		}
	}})
	informer.Run(informerStopper)
}

func checkInitialization() error {
	if client == nil {
		return fmt.Errorf("pipelines client is not initialized")
	}
	return nil
}
