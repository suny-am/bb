package api

type (
	Repositories struct {
		Values []Repository
	}
	Repository struct {
		Description string
		Name        string
		Size        int
		Language    string
		Project     Type
		Fork_Policy string
		Full_Name   string
		Is_Private  bool
		Has_Wiki    bool
		Owner       Owner
		Mainbranch  Type
		Readme      string

		Created_On string
		Updated_On string
	}
	Owner struct {
		Display_Name string
		Username     string
		UUID         string
		Type         string
	}
	Type struct {
		Name string
		Type string
	}
)
