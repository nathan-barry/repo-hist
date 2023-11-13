package handlers

type RawCommit struct {
	URL    string `json:"url"` // sha to commit files
	Commit struct {
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		}
		Message string `json:"message"`
	}
	Selected bool
}

type CommitData struct {
	Commit struct {
		Tree struct {
			URL string `json:"url"`
		}
	}
	Files []File `json:"files"`
}

type Dir struct {
	Tree []struct {
		Path      string `json:"path"`
		URL       string `json:"url"`
		Additions int
		Deletions int
		Status    string
	} `json:"tree"`
}

type Content struct {
	Content string `json:"content"`
}

type File struct {
	FileName    string `json:"filename"`
	Changes     int    `json:"changes"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	BlobURL     string `json:"blob_url"` // Link to github
	ContentsURL string `json:"contents_url"`
	RawURL      string `json:"raw_url"` // Contains entire file in resp.Body
	Patch       string `json:"patch"`
	SHA         string `json:"sha"`
	Status      string `json:"status"`
}

type ChangeData struct {
	Additions int
	Deletions int
	Status    string
}
