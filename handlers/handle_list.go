package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	CONTENTS_URL = "https://api.github.com/repos/github/gitignore/contents/"

	TYPE_FILE = "file"
	TYPE_DIR  = "dir"

	LIST_HEADER = "Available .gitignore files:\n\n"
	SEP         = "---\n"
)

type RepoContent []ContentEntry

type ContentEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func List(client Client) error {
	fileData, err := fetchList(client, CONTENTS_URL)
	if err != nil {
		return fmt.Errorf("could not fetch file list from %s: %w", CONTENTS_URL, err)
	}
	fileList := loadFiles(fileData)
	displayFileList(fileList)
	return nil
}

func loadFiles(content RepoContent) []string {
	fileList := []string{}
	for _, file := range content {
		if file.Type == TYPE_DIR {
			continue
		}
		if filepath.Ext(file.Name) != GIT_IGNORE {
			continue
		}
		fileList = append(fileList, file.Name)
	}
	return fileList
}

func fetchList(client Client, url string) (RepoContent, error) {
	var content RepoContent
	resp, err := client.Get(url)
	if err != nil {
		return content, err

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return content, fmt.Errorf("received non-OK HTTP status for %s: %s", url, resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return content, fmt.Errorf("failed to decode body: %w", err)
	}
	return content, nil

}

func displayFileList(files []string) {
	var currentFirst rune
	slices.SortFunc(files, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	fmt.Print(LIST_HEADER)
	for _, f := range files {
		// get first complete rune, not only first byte
		// but all files shoulde be ascii only anyway...
		firstLetter, _ := utf8.DecodeRuneInString(f)
		firstLetter = unicode.ToLower(firstLetter)
		if currentFirst == 0 {
			currentFirst = unicode.ToLower(firstLetter)
		}
		if firstLetter != currentFirst {
			fmt.Print(SEP)
			currentFirst = unicode.ToLower(firstLetter)
		}
		displayName := strings.TrimSuffix(f, filepath.Ext(f))
		fmt.Println(displayName)
	}
}
