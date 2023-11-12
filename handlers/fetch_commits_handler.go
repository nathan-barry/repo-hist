package handlers

import (
	"fmt"
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

func FetchCommitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinged -> FetchCommits")

	url := "https://api.github.com/repos/" + r.FormValue("repoURL") + "/commits"

	lastNum, _ := getLastPageNumber(url)

	rawCommits := fetchLastTenCommits(url, lastNum)

	t := template.Must(template.ParseFiles("./views/home/commits.html"))

	data := map[string]any{
		"RawCommits": rawCommits,
	}

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println("Template error:", err) // Log the error
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}
