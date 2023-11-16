package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchDirHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> Fetch Dir Handler")

	url := r.FormValue("url")

	// Grabs the commit metadata
	var commitData CommitData
	if err := getJSON(url, &commitData); err != nil {
		fmt.Println("getJSON error:", url)
		return
	}

	changedFiles, deletedFiles := processChangeFilesInfo(commitData)

	// Grabs all the files in the commit, adds commit metadata
	var dir Dir
	if err := getJSON(commitData.Commit.Tree.URL+"?recursive=1", &dir); err != nil {
		fmt.Println("getJSON error:", url)
		return
	}
	addChangedFileInfoToDir(changedFiles, &dir)
	addDeletedFilesToDir(deletedFiles, &dir)

	// Grabs content of first file
	fileURL, patch, path := getFirstChangedFileInfo(dir)

	var content Content
	if err := getJSON(fileURL, &content); err != nil {
		fmt.Println("getJSON error:", url)
		fmt.Println(err)
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(content.Content)

	// Template stuff
	t := template.Must(template.ParseFiles(
		"./views/home/dir.html",
		"./views/home/file.html",
	))

	data := map[string]any{
		"Tree": dir.Tree,
		// "InitialFetch": false,
		"File":  string(decoded),
		"Path":  path,
		"Patch": patch,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
