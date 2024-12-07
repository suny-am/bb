package api

type FileMetaResponse struct {
	Commit       Commit
	Links        Links
	Path         string
	Type         string
	Escaped_path string
	Mimetype     string
	Size         int
}
