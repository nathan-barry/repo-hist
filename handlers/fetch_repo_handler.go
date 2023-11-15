package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchRepoHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits", r.FormValue("repoURL"))
	lastNum, _ := getLastPageNumber(url)
	rawCommits := fetchLastCommits(url, lastNum)

	// Grabs the first commit metadata
	firstCommitURL := rawCommits[0].URL
	var commitData CommitData
	getJSON(firstCommitURL, &commitData, githubKey)
	changedFiles, deletedFiles := processChangeFilesInfo(commitData)

	// Grabs all the files in the commit, adds commit metadata
	var dir Dir
	getJSON(commitData.Commit.Tree.URL+"?recursive=1", &dir, githubKey)
	addChangedFileInfoToDir(changedFiles, &dir)
	addDeletedFilesToDir(deletedFiles, &dir)

	// Grabs content of first file
	firstFileURL := dir.Tree[0].URL
	var content Content
	getJSON(firstFileURL, &content, githubKey)
	decoded, _ := base64.StdEncoding.DecodeString(content.Content)
	path := dir.Tree[0].Path

	// Template stuff
	t := template.Must(template.ParseFiles(
		"./views/home/repo.html",
		"./views/home/file.html",
		"./views/home/dir.html",
	))

	data := map[string]any{
		"RawCommits":   rawCommits,
		"Tree":         dir.Tree,
		"InitialFetch": true,
		"File":         string(decoded),
		"Path":         path,
		"Patch":        commitData.Files[0].Patch,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err)
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
