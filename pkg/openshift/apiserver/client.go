package openshift

import (
	apiserverClientset "github.com/openshift/client-go/apiserver/clientset/versioned"
	"k8s.io/client-go/rest"
)

func NewClient() (*apiserverClientset.Clientset, error) {
	// Initially only allow handle same cluster, the solution should run inside a cluster
	// with SA allowing to create the required resources to handle events
	restconfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := apiserverClientset.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
