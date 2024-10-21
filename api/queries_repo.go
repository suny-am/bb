package api

type (
	Repositories struct {
		Size     int
		Page     int
		Previous string
		Next     string
		Values   []Repository
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
		Owner       User
		Mainbranch  Type
		Links       Links
		Readme      string

		Created_On string
		Updated_On string
	}
	Type struct {
		Name string
		Type string
	}
)
