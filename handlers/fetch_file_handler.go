package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> Fetch File Handler")

	url := r.FormValue("url")
	path := r.FormValue("path")
	patch := r.FormValue("patch")

	// Grabs content and changes of file
	var content Content
	if err := getJSON(url, &content); err != nil {
		fmt.Println("getJSON error:", url)
		return
	}
	decodedContent, _ := base64.StdEncoding.DecodeString(content.Content)
	decodedPatch, _ := base64.StdEncoding.DecodeString(patch)

	// Template Stuff
	t := template.Must(template.ParseFiles("./views/home/file.html"))

	data := map[string]any{
		"File":  string(decodedContent),
		"Path":  path,
		"Patch": string(decodedPatch),
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err)
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
