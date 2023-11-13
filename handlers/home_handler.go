package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var githubKey string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}
	githubKey = os.Getenv("GITHUB_AUTH")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Pinged -> Home")

	t := template.Must(template.ParseFiles("./views/partials/base.html", "./views/home/index.html"))

	data := map[string]any{
		"Title": "Repo History",
	}

	err := t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
