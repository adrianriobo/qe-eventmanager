package tekton

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianriobo/qe-eventmanager/pkg/util"
	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	informers "github.com/tektoncd/pipeline/pkg/client/informers/externalversions"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"knative.dev/pkg/apis"
)

type tektonClient struct {
	clientset       *clientset.Clientset
	namespace       string
	workspaces      []v1beta1.WorkspaceBinding
	defaultDuration v1.Duration
	consoleURL      string
}

type WorkspaceBinding struct {
	Name string
	PVC  string
}

var _client *tektonClient

func CreateClient(kubeconfig []byte, namespace string,
	workspaces []WorkspaceBinding, consoleURL string) (err error) {
	_client = &tektonClient{}
	_client.clientset, err = createClientset(kubeconfig)
	if len(namespace) > 0 {
		_client.namespace = namespace
	} else {
		err = fmt.Errorf("need to define a default namespace for tekton pipelineruns")
	}
	if len(workspaces) > 0 {
		fillWorkspaceBinding(workspaces)
	}
	// Move this to providers configuration file
	_client.defaultDuration = v1.Duration{Duration: 8 * time.Hour}
	_client.consoleURL = consoleURL
	if err == nil {
		logging.Debug("Tekton client initialized successfully")
	}
	return
}

func GetDefaultNamespace() string {
	return _client.namespace
}

func GetPipelinerunDashboardUrl(pipelinerunName string) string {
	return fmt.Sprintf("%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s",
		_client.consoleURL, _client.namespace, pipelinerunName)
}

func ApplyPipelinerun(spec *v1beta1.PipelineRun) (*v1beta1.PipelineRun, error) {
	if err := checkInitialization(); err != nil {
		return nil, err
	}
	spec.ObjectMeta.Namespace = _client.namespace
	spec.Spec.Workspaces = _client.workspaces
	spec.Spec.Timeout = &_client.defaultDuration
	return _client.clientset.TektonV1beta1().
		PipelineRuns(_client.namespace).
		Create(context.Background(), spec, v1.CreateOptions{})
}

// func(options *v1.ListOptions) {
// 	options.FieldSelector = fields.OneTermEqualSelector("metadata.name", "app-config").String()
// }
func AddInformer(pipelinerunName string, status chan *v1beta1.PipelineRunStatus, informerStopper chan struct{}) {
	if err := checkInitialization(); err != nil {
		logging.Error(err)
	}
	pipelineRunsInformer := informers.NewSharedInformerFactoryWithOptions(
		_client.clientset, 2*time.Minute,
		informers.WithNamespace(_client.namespace)).
		Tekton().V1beta1().PipelineRuns().Informer()
	pipelineRunsInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			pipelineRun, ok := obj.(*v1beta1.PipelineRun)
			if !ok {
				logging.Error("error formatting pipelinerun")
				return false
			}
			if pipelineRun.GetName() == pipelinerunName &&
				notifyStatus(pipelineRun) {
				return true
			}
			return false
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pipelineRun, _ := obj.(*v1beta1.PipelineRun)
				status <- &pipelineRun.Status
			},
		},
	})
	pipelineRunsInformer.Run(informerStopper)
}

func notifyStatus(pipelineRun *v1beta1.PipelineRun) bool {
	if util.IsEmpty(pipelineRun.Status) {
		return false
	}
	condition := pipelineRun.Status.GetCondition(apis.ConditionSucceeded)
	return (condition.Reason == string(v1beta1.PipelineRunReasonSuccessful) && waitForResults(pipelineRun)) ||
		condition.Reason == string(v1beta1.PipelineRunReasonFailed)
}

func waitForResults(pipelineRun *v1beta1.PipelineRun) bool {
	return len(pipelineRun.Status.PipelineRunStatusFields.PipelineSpec.Results) > 0 &&
		len(pipelineRun.Status.PipelineRunStatusFields.PipelineResults) > 0
}

func checkInitialization() error {
	if _client == nil {
		return fmt.Errorf("pipelines client is not initialized")
	}
	return nil
}

// Initially only allow handle same cluster, the solution should run inside a cluster
// with SA allowing to create the required resources to handle events
// var config
func createClientset(kubeconfig []byte) (client *clientset.Clientset, err error) {
	var config *rest.Config
	if len(kubeconfig) > 0 {
		config, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if config != nil {
		client, err = clientset.NewForConfig(config)
	}
	return
}

func fillWorkspaceBinding(workspacesInfo []WorkspaceBinding) {
	for _, workspaceInfo := range workspacesInfo {
		var workspace v1beta1.WorkspaceBinding = v1beta1.WorkspaceBinding{
			Name: workspaceInfo.Name,
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: workspaceInfo.PVC},
		}
		_client.workspaces = append(_client.workspaces, workspace)
	}
}
