package providers

type Providers struct {
	UMB    UMB    `yaml:"umb,omitempty"`
	Tekton Tekton `yaml:"tekton,omitempty"`
	Github Github `yaml:"github,omitempty"`
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
	ConsoleURL string      `yaml:"consoleURL"`
	Workspaces []Workspace `yaml:"workspaces"`
	Kubeconfig string      `yaml:"kubeconfig-data,omitempty"`
}

type Workspace struct {
	Name string `yaml:"name"`
	PVC  string `yaml:"pvc"`
}

type Github struct {
	Token string `yaml:"token"`
}
