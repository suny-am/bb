package api

type (
	Commit struct {
		Author  Author
		Date    string
		Hash    string
		Message string
		Parents []ParentCommit
	}
	Author struct {
		Raw  string
		Type string
		User User
	}
	CommitUser struct {
		Account_Id   string
		Display_Name string
		Nickname     string
		Type         string
		Uuid         string
	}
	ParentCommit struct {
		Hash string
		Type string
	}
)
