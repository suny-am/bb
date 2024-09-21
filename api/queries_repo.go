package api

type (
	Repositories struct {
		Values []Repository
	}
	Repository struct {
		Created_On  string
		Updated_On  string
		Description string
		Full_Name   string
		Is_Private  bool
	}
)
