package flows

type Flow struct {
	Name   string `yaml:"name"`
	Input  Input  `yaml:"input"`
	Action Action `yaml:"action"`
}

type Input struct {
	UMB UMBInput `yaml:"umb,omitempty"`
}

type UMBInput struct {
	Topic   string   `yaml:"topic"`
	Filters []string `yaml:"filters"`
}

type Action struct {
	TektonPipeline TektonPipelineAction `yaml:"tektonPipeline,omitempty"`
	Forward        ForwardAction        `yaml:"forward,omitempty"`
}

type TektonPipelineAction struct {
	Name    string          `yaml:"name"`
	Params  []NameValuePair `yaml:"params"`
	Success Success         `yaml:"success,omitempty"`
	Error   Error           `yaml:"error,omitempty"`
}

type ForwardAction struct {
	Type string `yaml:"type"`
}

type NameValuePair struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}

type Success struct {
	UMB UMBEvent `yaml:"umb"`
}

type Error struct {
	UMB UMBEvent `yaml:"umb"`
}

type UMBEvent struct {
	Topic       string          `yaml:"topic"`
	EventSchema string          `yaml:"eventSchema"`
	EventFields []NameValuePair `yaml:"eventFields"`
}
