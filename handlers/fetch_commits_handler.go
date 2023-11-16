package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func FetchCommitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> Fetch Commits Handler")

	var url string
	var rawCommits []*RawCommit
	var numCurPage int
	var numTotalPages int

	initialSearch := r.FormValue("initial_search")

	if initialSearch == "true" {
		fmt.Println("INITIAL SEARCH")
		url = fmt.Sprintf("https://api.github.com/repos/%s/commits", r.FormValue("repoURL"))
		lastNum, err := getLastPageNumber(url)
		numCurPage = 1
		numTotalPages = lastNum
		if err != nil {
			fmt.Println("LastPageNumber error:", url)
			return
		}
		rawCommits = fetchCommits(url, lastNum)
	} else {
		fmt.Println("NOT INITIAL SEARCH")
		url = r.FormValue("url")
		action := r.FormValue("action")

		numCurPage, _ = strconv.Atoi(r.FormValue("num_cur_page"))
		numTotalPages, _ = strconv.Atoi(r.FormValue("num_total_pages"))

		if action == "prev" {
			numCurPage -= 1
		} else if action == "next" {
			numCurPage += 1
		} else {
			// do something
		}

		rawCommits = fetchCommits(url, numTotalPages-numCurPage+1)
	}

	// Grabs the first commit metadata
	firstCommitURL := rawCommits[0].URL
	var commitData CommitData
	if err := getJSON(firstCommitURL, &commitData); err != nil {
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
	firstFileURL := dir.Tree[0].URL
	var content Content
	if err := getJSON(firstFileURL, &content); err != nil {
		fmt.Println("getJSON error:", url)
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(content.Content)
	path := dir.Tree[0].Path

	// Template stuff
	t := template.Must(template.ParseFiles(
		"./views/home/repo.html",
		"./views/home/file.html",
		"./views/home/dir.html",
		"./views/home/commit_list.html",
	))

	data := map[string]any{
		"RawCommits":    rawCommits,
		"Tree":          dir.Tree,
		"InitialFetch":  true,
		"File":          string(decoded),
		"Path":          path,
		"Patch":         commitData.Files[0].Patch,
		"NumTotalPages": numTotalPages,
		"NumCurPage":    numCurPage,
		"URL":           url,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err)
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
