package api

type (
	PullRequests struct {
		Values []PullRequest
	}
	PullRequest struct {
		Author              Author
		Close_Source_Branch bool
		Comment_Count       int
		Description         string
		Destination         Destination
		Id                  int
		Merge_Commit        Commit
		Links               Links
		Reason              string
		Source              Source
		State               string
		Task_Count          int
		Title               string
		Type                string
		Updated_On          string
	}
	Links struct {
		Self     map[string]string
		Html     map[string]string
		Commits  map[string]string
		Approve  map[string]string
		Diff     map[string]string
		DiffStat map[string]string
		Comments map[string]string
		Activity map[string]string
		Merge    map[string]string
		Decline  map[string]string
	}
	Source struct {
		Branch     Branch
		Commit     Commit
		Repository PrRepository
	}
	Destination struct {
		Branch       Branch
		Commit       Commit
		PrRepository PrRepository
	}
	PrCommit struct {
		Hash string
		Type string
	}
	PrRepository struct {
		Full_Name string
		Name      string
		Type      string
		Uuid      string
	}
	Branch struct {
		Name string
	}
	PrAuthor struct {
		Account_Id   string
		Display_Name string
		Nickname     string
		Type         string
		Uuid         string
	}
)
