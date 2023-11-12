package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func FetchDirHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinged -> FetchDir")

	url := r.FormValue("url")

	fmt.Println("THIS IS THE DIR URL ->", url)

	var commitFiles Files
	getJSON(url, &commitFiles, githubKey)

	var dirURL DirURL
	getJSON(url, &dirURL, githubKey)

	printJSON("DirURL", dirURL)

	var dir Dir
	getJSON(dirURL.Commit.Tree.URL+"?recursive=1", &dir, githubKey)

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
