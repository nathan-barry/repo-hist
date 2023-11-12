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
}

type DirURL struct {
	Commit struct {
		Tree struct {
			URL string `json:"url"`
		}
	}
}

type Dir struct {
	SHA  string `json:"sha"`
	URL  string `json:"url"`
	Tree []struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	} `json:"tree"`
}

type Content struct {
	Content string `json:"content"`
}

type Files struct {
	Files []File `json:"files"`
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
