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

	t := template.Must(template.ParseFiles("./views/home/repo.html"))

	// Grab first commit dir
	firstCommitURL := rawCommits[0].URL

	var dirURL DirURL
	getJSON(firstCommitURL, &dirURL, githubKey)

	var dir Dir
	getJSON(dirURL.Commit.Tree.URL+"?recursive=1", &dir, githubKey)

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

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
