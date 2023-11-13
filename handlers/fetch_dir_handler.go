package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func FetchDirHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	var commitData CommitData
	getJSON(url, &commitData, githubKey)

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

	t := template.Must(template.ParseFiles("./views/home/dir.html"))

	data := map[string]any{
		"Tree": dir.Tree,
	}

	fmt.Println(dir.Tree)

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
