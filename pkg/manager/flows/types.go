package flows

type Flow struct {
	Name   string `yaml:"name"`
	Input  Input  `yaml:"input"`
	Action Action `yaml:"action"`
}

type Input struct {
	UmbInput *UmbInput `yaml:"umb,omitempty"`
}

type UmbInput struct {
	Topic   string           `yaml:"topic"`
	Filters []UmbInputFilter `yaml:"filters"`
}

type UmbInputFilter struct {
	JsonPath string `yaml:"jsonpath"`
}

type Action struct {
	TektonPipelineAction *TektonPipelineAction `yaml:"tektonPipeline,omitempty"`
	Success              Success               `yaml:"success"`
	Error                Error                 `yaml:"error"`
}

type TektonPipelineAction struct {
	PipelineName   string                `yaml:"name"`
	PipelineParams []TektonPipelineParam `yaml:"params"`
}

type TektonPipelineParam struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value,omitempty"`
	JsonPath string `yaml:"jsonpath,omitempty"`
}

type Success struct {
	UMBEvent UMBEvent `yaml:"umb"`
}

type Error struct {
	UMBEvent UMBEvent `yaml:"umb"`
}

type UMBEvent struct {
	Topic  string          `yaml:"topic"`
	Schema string          `yaml:"eventSchema"`
	Fields []UMBEventField `yaml:"eventFields"`
}

type UMBEventField struct {
	Name                  string `yaml:"name"`
	Value                 string `yaml:"value,omitempty"`
	PipelineResultName    string `yaml:"pipelineResultName,omitempty"`
	PipelineParameterName string `yaml:"pipelineParameterName,omitempty"`
}
