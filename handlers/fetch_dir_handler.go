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

	t := template.Must(template.ParseFiles("./views/home/dir.html"))

	data := map[string]any{
		"FileArray": commitFiles.Files,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
