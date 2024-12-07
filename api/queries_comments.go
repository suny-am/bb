package api

type Comments struct {
	Values []Comment
}

type Comment struct {
	Links      Links
	User       User
	Content    Content
	Created_On string
	Updated_On string
	Inline     Inline
	Id         int
	Deleted    bool
	Pending    bool
}

type Inline struct {
	Path string
	From int
	To   int
}

type Content struct {
	Type string
	Raw  string
	Html string
}
