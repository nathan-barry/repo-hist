package handlers

import (
	"encoding/base64"
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
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Extract the Link header. If "" then only 1 page
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

func fetchCommits(url string, pageNum int) []*RawCommit {
	pageURL := fmt.Sprintf("%s?page=%v", url, pageNum)
	fmt.Println("Fetching commit from:", pageURL)

	var commits []*RawCommit
	getJSON(pageURL, &commits)

	// reverse list
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}
	return commits
}

func getFirstChangedFileInfo(dir Dir) (fileURL, patch, path string) {
	for i := 0; i < len(dir.Tree); i++ {
		if dir.Tree[i].Patch != "" {
			fileURL = dir.Tree[i].URL
			decodedPatch, _ := base64.StdEncoding.DecodeString(dir.Tree[i].Patch)
			patch = string(decodedPatch)
			path = dir.Tree[i].Path
			return
		}
	}
	return
}

func processChangeFilesInfo(commitData CommitData) (map[string]ChangeData, []string) {
	changedFiles := make(map[string]ChangeData, len(commitData.Files))
	deletedFiles := []string{}
	for _, f := range commitData.Files {
		if f.Status == "removed" {
			deletedFiles = append(deletedFiles, f.FileName)
		} else {
			changedFiles[f.FileName] = ChangeData{
				Additions: f.Additions,
				Deletions: f.Deletions,
				Patch:     base64.StdEncoding.EncodeToString([]byte(f.Patch)),
			}
		}
	}

	return changedFiles, deletedFiles
}

func addDeletedFilesToDir(deletedFiles []string, dir *Dir) {
	for _, f := range deletedFiles {
		tree := struct {
			Path      string `json:"path"`
			URL       string `json:"url"`
			Additions int
			Deletions int
			Patch     string
		}{
			Path: f,
			URL:  "deleted",
		}
		dir.Tree = append(dir.Tree, tree)
	}
}

func addChangedFileInfoToDir(changedFiles map[string]ChangeData, dir *Dir) {
	for i := 0; i < len(dir.Tree); i++ {
		if cd, ok := changedFiles[dir.Tree[i].Path]; ok {
			dir.Tree[i].Deletions = cd.Deletions
			dir.Tree[i].Additions = cd.Additions
			dir.Tree[i].Patch = cd.Patch
		}
	}
}

func getBody(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// If key isn't empty, add it to header
	if githubKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", githubKey))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getJSON(url string, data any) error {
	body, err := getBody(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

// Helper to print after getJSON for debugging
func printJSON(data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("PrettyJSON Error", err)
		log.Fatal(err)
	}

	fmt.Println("PRINT_JSON: ", string(prettyJSON))
}
