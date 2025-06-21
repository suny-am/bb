package api

type (
	BaseOptions struct {
		Credentials string
		Workspace   string
		Repository  string
	}

	CodeSearchOptions struct {
		BaseOptions
		QueryParameters

		IncludeSource bool
	}

	PipelineListOptions struct {
		BaseOptions
		QueryParameters
	}

	PullrequestViewOptions struct {
		BaseOptions
		QueryParameters

		Pullrequest string
	}

	PullrequestListOptions struct {
		BaseOptions
		QueryParameters

		Title     string
		Creator   string
		State     string
		Approvals int
	}

	RepositoryViewOptions struct {
		BaseOptions
	}

	RepositoryListOptions struct {
		BaseOptions
		QueryParameters

		Name string
	}

	ForkListptions struct {
		BaseOptions
	}
)
