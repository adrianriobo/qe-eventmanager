package pipelines

import (
	"time"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	crcNamespace string = "codeready-container"
)

var (
	crcWorkspace v1beta1.WorkspaceBinding = v1beta1.WorkspaceBinding{
		Name: "pipelines-data",
		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			ClaimName: "pipelines-data"},
	}

	defaultTimeout v1.Duration = v1.Duration{
		Duration: 5 * time.Hour,
	}
)
