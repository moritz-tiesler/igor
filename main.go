package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const (
	contentsUrl = "https://api.github.com/repos/github/gitignore/contentsss/"
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
		fmt.Println(fileList)
	}
}

const (
	TYPE_FILE = "file"
	TYPE_DIR  = "dir"
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
