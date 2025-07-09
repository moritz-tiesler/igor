package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	contentsUrl = "https://api.github.com/repos/github/gitignore/contents/"
	rawPrefix   = "https://raw.githubusercontent.com/github/gitignore/main/"
)

var doList bool

func init() {
	const (
		listDefault bool   = false
		usage       string = "list available .gitignore files"
	)
	flag.BoolVar(&doList, "list", listDefault, usage)
}

func main() {
	flag.Parse()
	if doList {
		fileData, err := fetchList(contentsUrl)
		if err != nil {
			fmt.Printf("could not fetch file list form %s\n", contentsUrl)
			os.Exit(1)
		}
		fileList := loadFiles(fileData)
		displayFileList(fileList)
	}
}

const (
	TYPE_FILE = "file"
	TYPE_DIR  = "dir"
	EXT       = ".gitignore"
)

type Content []ContentEntry

type ContentEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func loadFiles(content Content) []string {
	fileList := []string{}
	for _, file := range content {
		if file.Type == TYPE_DIR {
			continue
		}
		if filepath.Ext(file.Name) != EXT {
			continue
		}
		fileList = append(fileList, file.Name)
	}
	return fileList
}

func fetchList(url string) (Content, error) {
	var content Content
	data, err := http.Get(url)
	if err != nil {
		return content, err
	}
	err = json.NewDecoder(data.Body).Decode(&content)
	if err != nil {
		return content, err
	}
	return content, nil

}

const LIST_HEADER = "Available .gitignore files:\n\n"
const SEP = "---\n"

// assumes the files are sorted alphabetically
func displayFileList(files []string) {
	var currentFirst rune
	fmt.Print(LIST_HEADER)
	for _, f := range files {
		// get first complete rune, not only first byte
		firstLetter, _ := utf8.DecodeRuneInString(f)
		if currentFirst == 0 {
			currentFirst = firstLetter
		}
		if firstLetter != currentFirst {
			fmt.Print(SEP)
			currentFirst = firstLetter
		}
		displayName := strings.TrimSuffix(f, filepath.Ext(f))
		fmt.Println(displayName)
	}
}
