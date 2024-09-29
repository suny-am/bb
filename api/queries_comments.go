package api

type Comments struct {
	Values []Comment
}

type Comment struct {
	Id         int
	Created_On string
	Updated_On string
	Content    Content
	User       User
	Deleted    bool
	Inline     Inline
	Pending    bool
	Links      Links
}

type Inline struct {
	From int
	To   int
	Path string
}

type Content struct {
	Type string
	Raw  string
	Html string
}
