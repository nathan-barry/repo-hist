package handlers

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func FetchFileHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	path := r.FormValue("path")

	fmt.Println("PATH", path)

	var content Content
	getJSON(url, &content, githubKey)

	decoded, _ := base64.StdEncoding.DecodeString(content.Content)

	t := template.Must(template.ParseFiles("./views/home/file.html"))

	data := map[string]any{
		"File": string(decoded),
		"Path": path,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
