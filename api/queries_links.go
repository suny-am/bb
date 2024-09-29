package api

type Links struct {
	Self     Link
	Html     Link
	Commits  Link
	Approve  Link
	Diff     Link
	Diffstat Link
	Comments Link
	Activity Link
	Merge    Link
	Decline  Link
}

type Link struct {
	Href string
	Name string
}
