package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var githubKey string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}
	githubKey = os.Getenv("GITHUB_AUTH")
}

func FetchRepoHandler(w http.ResponseWriter, r *http.Request) {
	// Grab repo info
	url := "https://api.github.com/repos/" + r.FormValue("repoURL") + "/commits"

	lastNum, _ := getLastPageNumber(url)
	rawCommits := fetchLastTenCommits(url, lastNum)

	rawCommits[0].Selected = true

	// Grab first commit dir
	firstCommitURL := rawCommits[0].URL

	var commitData CommitData
	getJSON(firstCommitURL, &commitData, githubKey)

	changedFiles := make(map[string]ChangeData, len(commitData.Files))
	for _, f := range commitData.Files {
		changedFiles[f.FileName] = ChangeData{
			Additions: f.Additions,
			Deletions: f.Deletions,
			Status:    f.Status,
		}
	}

	var dir Dir
	getJSON(commitData.Commit.Tree.URL+"?recursive=1", &dir, githubKey)

	for i := 0; i < len(dir.Tree); i++ {
		if cd, ok := changedFiles[dir.Tree[i].Path]; ok {
			dir.Tree[i].Deletions = cd.Deletions
			dir.Tree[i].Additions = cd.Additions
			dir.Tree[i].Status = cd.Status
		}
	}

	// Grab first file
	firstFileURL := dir.Tree[0].URL
	path := dir.Tree[0].Path

	var content Content
	getJSON(firstFileURL, &content, githubKey)

	decoded, _ := base64.StdEncoding.DecodeString(content.Content)

	// Template stuff below
	data := map[string]any{
		"RawCommits": rawCommits,
		"Tree":       dir.Tree,
		"File":       string(decoded),
		"Path":       path,
	}

	t := template.Must(template.ParseFiles("./views/home/repo.html"))

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
