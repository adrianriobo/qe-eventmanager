package pipelines

import (
	"time"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace          string = "codeready-container"
	DashboardBaseUrl   string = "https://tekton-dashboard-openshift-pipelines.apps.ocp4.prod.psi.redhat.com"
	DashboardUrlFormat string = "%s/#/namespaces/%s/pipelineruns/%s"
)

var (
	Workspace v1beta1.WorkspaceBinding = v1beta1.WorkspaceBinding{
		Name: "pipelines-data",
		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			ClaimName: "pipelines-data"},
	}

	DefaultTimeout v1.Duration = v1.Duration{
		Duration: 8 * time.Hour,
	}
)
