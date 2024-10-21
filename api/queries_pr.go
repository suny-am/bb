package api

type (
	Pullrequests struct {
		Size     int
		Page     int
		Previous string
		Next     string
		Values   []Pullrequest
	}
	Pullrequest struct {
		Id                  int
		Summary             TextElem
		Links               Links
		Title               string
		Rendered            Rendered
		State               string
		Author              User
		Comment_Count       int
		Task_Count          int
		Merge_Commit        Commit
		Close_Source_Branch bool
		Closed_By           Type
		Source              Source
		Reason              string
		Created_On          string
		Updated_On          string
		Reviewers           []User
		Participants        []Participant
		Comments            Comments
	}
	Participant struct {
		Type     string
		User     User
		Role     string
		Approved bool
	}
	Source struct {
		Repository Repository
		Branch     Branch
		Commit     Commit
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
