package api

type (
	Repositories struct {
		Previous string
		Next     string
		Values   []Repository
		Size     int
		Page     int
	}
	Repository struct {
		Links       Links
		Owner       User
		Project     Type
		Mainbranch  Type
		Language    string
		Fork_Policy string
		Full_Name   string
		Description string
		Name        string
		Readme      string
		Created_On  string
		Updated_On  string
		Size        int
		Is_Private  bool
		Has_Wiki    bool
	}
	Type struct {
		Name string
		Type string
		Key  string
	}
)
