package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func FetchDirHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	// Grabs the commit metadata
	var commitData CommitData
	getJSON(url, &commitData, githubKey)
	changedFiles, deletedFiles := processChangeFilesInfo(commitData)

	// Grabs all the files in the commit, adds commit metadata
	var dir Dir
	getJSON(commitData.Commit.Tree.URL+"?recursive=1", &dir, githubKey)
	addChangedFileInfoToDir(changedFiles, &dir)
	addDeletedFilesToDir(deletedFiles, &dir)

	// Template stuff
	t := template.Must(template.ParseFiles("./views/home/dir.html"))

	data := map[string]any{
		"Tree":         dir.Tree,
		"InitialFetch": false,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
