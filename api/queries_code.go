package api

type CodeSearchResponse struct {
	Size              int
	Page              int
	PageLen           int
	Next              string
	Query_Substituted bool
	Values            []CodeItem
}

type CodeItem struct {
	Type                string
	Content_match_count int
	Content_matches     []ContentMatch
	Patch_matches       []PathMatch
	File                File
}

type ContentMatch struct {
	Lines []Line
}

type Line struct {
	Line     int
	Segments []Segment
}

type Segment struct {
	Text  string
	Match bool
}

type PathMatch struct {
}

type File struct {
	Path   string
	Type   string
	Links  Links
	Source string
}
