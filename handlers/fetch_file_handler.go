package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func FetchFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinged -> FetchCode")

	url := r.FormValue("raw_url")
	fileName := r.FormValue("file_name")
	body := getBody(url, githubKey)

	t := template.Must(template.ParseFiles("./views/home/file.html"))

	data := map[string]any{
		"Code":     string(body),
		"FileName": fileName,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
