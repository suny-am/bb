package api

type Pipelines struct {
	Previous string
	Next     string
	Values   []Pipeline
	Size     int
	Page     int
}

type Pipeline struct {
	Target                Branch
	Links                 Links
	State                 PipelineState
	Creator               User
	Type                  string
	UUID                  string
	Created_On            string
	Completed_On          string
	Configuration_Sources []Configuration_Source
	Repository            Repository
	Build_Number          int
	Build_Seconds_Used    int
}

type PipelineState struct {
	Name   string
	Type   string
	Stage  PipelineStage
	Result PipelineResult
}

type PipelineStage struct {
	Name string
	Type string
}

type PipelineResult struct {
	Name  string
	Type  string
	Error PipelineError
}

type PipelineError struct {
	Type    string
	Key     string
	Message string
}

type Configuration_Source struct {
	Source Branch
	URI    string
}
