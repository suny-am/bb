package api

type FileMetaResponse struct {
	Path         string
	Commit       Commit
	Type         string
	Escaped_path string
	Size         int
	Mimetype     string
	Links        Links
}
