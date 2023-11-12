package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

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

func fetchLastTenCommits(url string, lastNum int) []*RawCommit {
	lastURL := fmt.Sprintf("%s?page=%v", url, lastNum)

	var commits []*RawCommit
	getJSON(lastURL, &commits, githubKey)

	// reverse
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	return commits
}

// Takes in the URL and an optional github api key
func getBody(url string, key string) []byte {
	fmt.Println("\nGetting from URL:", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// If key isn't empty, add it to header
	if key != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", key))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("STATUS CODE:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func getJSON(url string, data any, key string) {
	body := getBody(url, key)
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Get JSON Fucked up", err)
		log.Fatal(err)
	}
}

func prettyJSON(data any) []byte {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("PrettyJSON Fucked up")
		log.Fatal(err)
	}
	return prettyJSON
}

func printJSON(name string, data any) {
	fmt.Println("PRINT_JSON: " + name + ":\n" + string(prettyJSON(data)))
}
