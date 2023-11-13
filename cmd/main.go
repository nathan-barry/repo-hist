package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nathan-barry/repo-hist/handlers"
)

func main() {
	fmt.Println("Starting server...")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/fetch-repo", handlers.FetchRepoHandler)
	http.HandleFunc("/fetch-dir", handlers.FetchDirHandler)
	http.HandleFunc("/fetch-file", handlers.FetchFileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func main() {
// 	// Example usage
// 	repoName := "nathan-barry/pretty-commit" // replace with the desired repository
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/commits", repoName)
// 	lastPage, err := getLastPageNumber(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Last page number: %d\n", lastPage)

// 	lastURL := fmt.Sprintf("%s?page=%v", url, lastPage)

// 	var commits []*Commit
// 	getJSON(lastURL, &commits)

// 	count := 0
// 	for _, c := range commits {
// 		fmt.Println(c.URL)
// 		count += 1
// 	}

// 	var commitData any
// 	getJSON(commits[0].URL, &commitData)
// 	prettyJSON, _ := json.MarshalIndent(commitData, "", "    ")
// 	fmt.Println(string(prettyJSON))

// }
