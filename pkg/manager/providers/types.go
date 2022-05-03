package providers

type Providers struct {
	UMB    UMB    `yaml:"umb"`
	Tekton Tekton `yaml:"tekton"`
}

type UMB struct {
	ConsumerID           string `yaml:"consumerID"`
	Driver               string `yaml:"driver"`
	Brokers              string `yaml:"brokers"`
	UserCertificate      string `yaml:"userCertificate"`
	UserKey              string `yaml:"userKey"`
	CertificateAuthority string `yaml:"certificateAuthority"`
}

type Tekton struct {
	Namespace  string      `yaml:"namespace"`
	Workspaces []Workspace `yaml:"workspaces"`
	Kubeconfig string      `yaml:"kubeconfig-data,omitempty"`
}

type Workspace struct {
	Name string `yaml:"name"`
	PVC  string `yaml:"pvc"`
}
