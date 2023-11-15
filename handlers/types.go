package handlers

type RawCommit struct {
	URL    string `json:"url"`
	Commit struct {
		Committer struct {
			Name string `json:"name"`
			Date string `json:"date"`
		}
		Message string `json:"message"`
	}
}

type CommitData struct {
	Commit struct {
		Tree struct {
			URL string `json:"url"`
		}
	}
	Files []File `json:"files"`
}

type File struct {
	FileName    string `json:"filename"`
	Changes     int    `json:"changes"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	BlobURL     string `json:"blob_url"`
	ContentsURL string `json:"contents_url"`
	RawURL      string `json:"raw_url"`
	Patch       string `json:"patch"`
	SHA         string `json:"sha"`
	Status      string `json:"status"`
}

type ChangeData struct {
	Additions int
	Deletions int
	Patch     string
}

type Dir struct {
	Tree []struct {
		Path      string `json:"path"`
		URL       string `json:"url"`
		Additions int
		Deletions int
		Patch     string
	} `json:"tree"`
}

type Content struct {
	Content string `json:"content"`
}
