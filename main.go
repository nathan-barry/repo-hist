package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	// Example usage
	repoName := "nathan-barry/pretty-commit" // replace with the desired repository
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits", repoName)
	lastPage, err := getLastPageNumber(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Last page number: %d\n", lastPage)

	lastURL := fmt.Sprintf("%s?page=%v", url, lastPage)

	var commits []*Commit
	getJSON(lastURL, &commits)

	count := 0
	for _, c := range commits {
		fmt.Println(c.URL)
		count += 1
	}

	var commitData any
	getJSON(commits[0].URL, &commitData)
	prettyJSON, _ := json.MarshalIndent(commitData, "", "    ")
	fmt.Println(string(prettyJSON))

}

// getLastPageNumber fetches the first page of commits for the specified repository
// and parses the `Link` header to find the last page number.
func getLastPageNumber(url string) (int, error) {
	// Perform a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Extract the Link header
	linkHeader := resp.Header.Get("Link")
	if linkHeader == "" {
		return 1, nil
	}

	// Regular expression to find the URL marked as "last"
	re := regexp.MustCompile(`<https://api\.github\.com/.*\?page=(\d+)>; rel="last"`)
	matches := re.FindStringSubmatch(linkHeader)
	if len(matches) != 2 {
		return 0, fmt.Errorf("could not parse last page number from Link header")
	}

	// Convert the page number to an integer
	lastPage, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid page number format: %s", err)
	}

	return lastPage, nil
}

func fetchLastTenCommits(url string, lastNum int) []*Commit {
	lastURL := fmt.Sprintf("%s?page=%v", url, lastNum)

	var commits []*Commit
	getJSON(lastURL, &commits)

	for _, c := range commits {
		fmt.Println(c.URL)
	}

	return commits
}

type Commit struct {
	URL string `json:"url"`
}

func getBody(url string) []byte {
	fmt.Println("Getting from URL:", url)
	// Get response from URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read body of response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// // Print the entire header
	// for name, values := range resp.Header {
	// 	// Each value is a slice, in case the header has multiple values
	// 	for _, value := range values {
	// 		fmt.Printf("%s: %s\n", name, value)
	// 	}
	// }

	return body
}

func getJSON(url string, data any) {
	body := getBody(url)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
}
