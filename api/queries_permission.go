package api

type (
	Permissions struct {
		Values []Permission
	}

	Permission struct {
		Permission string
		Type       string
		User       PermissionUser
		Repository PermissionRepository
	}

	PermissionUser struct {
		AccountId    string
		Display_Name string
		Nickname     string
		Type         string
		Uuid         string
	}

	PermissionRepository struct {
		Full_Name string
		Name      string
		DataType  string
		Uuid      string
		Type      string
	}
)
