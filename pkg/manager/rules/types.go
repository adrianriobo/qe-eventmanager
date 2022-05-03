package rules

type Rule struct {
	Name    string  `yaml:"name"`
	Input   Input   `yaml:"input"`
	Action  Action  `yaml:"action"`
	Success Success `yaml:"success"`
	Error   Error   `yaml:"error"`
}

type Input struct {
	UmbInput UmbInput `yaml:"umb,omitempty"`
}

type UmbInput struct {
	Topic   string   `yaml:"topic"`
	Filters []string `yaml:"filters"`
}

type Action struct {
	TektonPipelineAction TektonPipelineAction `yaml:"tektonPipeline,omitempty"`
}

type TektonPipelineAction struct {
	PipelineName   string  `yaml:"name"`
	PipelineParams []Tuple `yaml:"params"`
	Results        []Tuple `yaml:"results"`
}

type Success struct {
	UMBEvent UMBEvent `yaml:"umb"`
}

type Error struct {
	UMBEvent UMBEvent `yaml:"umb"`
}

type UMBEvent struct {
	Topic  string  `yaml:"topic"`
	Schema string  `yaml:"eventSchema"`
	Fields []Tuple `yaml:"eventFields"`
}

type Tuple struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
