package api

type CodeSearchResponse struct {
	Next              string
	Values            []CodeItem
	Size              int
	Page              int
	PageLen           int
	Query_Substituted bool
}

type CodeItem struct {
	File                File
	Type                string
	Content_matches     []ContentMatch
	Patch_matches       []PathMatch
	Content_match_count int
}

type ContentMatch struct {
	Lines []Line
}

type Line struct {
	Segments []Segment
	Line     int
}

type Segment struct {
	Text  string
	Match bool
}

type PathMatch struct{}

type File struct {
	Path   string
	Type   string
	Links  Links
	Source string
}
