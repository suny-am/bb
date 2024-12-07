package api

type (
	Pullrequests struct {
		Previous string
		Next     string
		Values   []Pullrequest
		Size     int
		Page     int
	}
	Pullrequest struct {
		Merge_Commit        Commit
		Links               Links
		Rendered            Rendered
		Author              User
		Summary             TextElem
		Closed_By           Type
		State               string
		Updated_On          string
		Title               string
		Reason              string
		Created_On          string
		Comments            Comments
		Participants        []Participant
		Reviewers           []User
		Source              Source
		Comment_Count       int
		Task_Count          int
		Id                  int
		Close_Source_Branch bool
	}
	Participant struct {
		Type     string
		User     User
		Role     string
		Approved bool
	}
	Source struct {
		Branch     Branch
		Commit     Commit
		Repository Repository
	}
	Branch struct {
		Name  string
		Links Links
	}
	Rendered struct {
		Title       TextElem
		Description TextElem
		Reason      TextElem
	}
	TextElem struct {
		Raw    string
		Markup string
		Html   string
	}
)
