package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchRepoHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("FetchRepoHandler...")

	// Grab repo info
	url := "https://api.github.com/repos/" + r.FormValue("repoURL") + "/commits"

	lastNum, _ := getLastPageNumber(url)
	rawCommits := fetchLastTenCommits(url, lastNum)

	// Grab first commit dir
	firstCommitURL := rawCommits[0].URL

	var commitData CommitData
	getJSON(firstCommitURL, &commitData, githubKey)

	changedFiles := make(map[string]ChangeData, len(commitData.Files))
	deletedFiles := []string{}
	for _, f := range commitData.Files {
		if f.Status == "removed" {
			deletedFiles = append(deletedFiles, f.FileName)
		} else {
			changedFiles[f.FileName] = ChangeData{
				Additions: f.Additions,
				Deletions: f.Deletions,
				Patch:     base64.StdEncoding.EncodeToString([]byte(f.Patch)),
			}
		}
	}

	var dir Dir
	getJSON(commitData.Commit.Tree.URL+"?recursive=1", &dir, githubKey)

	for i := 0; i < len(dir.Tree); i++ {
		if cd, ok := changedFiles[dir.Tree[i].Path]; ok {
			dir.Tree[i].Deletions = cd.Deletions
			dir.Tree[i].Additions = cd.Additions
			dir.Tree[i].Patch = cd.Patch
		}
	}

	for _, f := range deletedFiles {
		tree := struct {
			Path      string `json:"path"`
			URL       string `json:"url"`
			Additions int
			Deletions int
			Patch     string
		}{
			Path: f,
			URL:  "deleted",
		}
		dir.Tree = append(dir.Tree, tree)
	}

	// Grab first file
	firstFileURL := dir.Tree[0].URL
	path := dir.Tree[0].Path

	var content Content
	getJSON(firstFileURL, &content, githubKey)

	decoded, _ := base64.StdEncoding.DecodeString(content.Content)

	// Template stuff below
	data := map[string]any{
		"RawCommits":   rawCommits,
		"Tree":         dir.Tree,
		"InitialFetch": true,
		"File":         string(decoded),
		"Path":         path,
		"Patch":        commitData.Files[0].Patch,
	}

	t := template.Must(template.ParseFiles(
		"./views/home/repo.html",
		"./views/home/file.html",
		"./views/home/dir.html",
	))

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
