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
	http.HandleFunc("/fetch-dir", handlers.FetchDirHandler)
	http.HandleFunc("/fetch-file", handlers.FetchFileHandler)
	http.HandleFunc("/fetch-commits", handlers.FetchCommitsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
