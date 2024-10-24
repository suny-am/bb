package api

type Pipelines struct {
	Size     int
	Page     int
	Previous string
	Next     string
	Values   []Pipeline
}

type Pipeline struct {
	Type                  string
	UUID                  string
	Build_Number          int
	Creator               User
	Repository            Repository
	Target                Branch
	State                 PipelineState
	Created_On            string
	Completed_On          string
	Build_Seconds_Used    int
	Configuration_Sources []Configuration_Source
	Links                 Links
}

type PipelineState struct {
	Name   string
	Type   string
	Result PipelineResult
}

type PipelineResult struct {
	Name string
	Type string
}

type Configuration_Source struct {
	Source Branch
	URI    string
}
