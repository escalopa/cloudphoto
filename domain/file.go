package domain

type File struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Folder struct {
	Key     string    `json:"key"`
	Name    string    `json:"name"`
	Files   []*File   `json:"files,omitempty"`
	Folders []*Folder `json:"folders,omitempty"`
}
