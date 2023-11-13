package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FetchFileHandler...")

	url := r.FormValue("url")
	path := r.FormValue("path")
	patch := r.FormValue("patch")

	fmt.Println("PATH", path)

	var content Content
	getJSON(url, &content, githubKey)

	decodedContent, _ := base64.StdEncoding.DecodeString(content.Content)
	decodedPatch, _ := base64.StdEncoding.DecodeString(patch)

	t := template.Must(template.ParseFiles("./views/home/file.html"))

	data := map[string]any{
		"File":  string(decodedContent),
		"Path":  path,
		"Patch": string(decodedPatch),
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
