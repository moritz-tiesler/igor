package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	CONTENTS_URL = "https://api.github.com/repos/github/gitignore/contents/"
	RAWP_REFIX   = "https://raw.githubusercontent.com/github/gitignore/main/"
)

var doList bool

func init() {
	const (
		listDefault bool   = false
		listUsage   string = "list available .gitignore files"
	)
	flag.BoolVar(&doList, "list", listDefault, listUsage)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  This tool copies a .gitignore file for a specified language into the current directory.\n")
		fmt.Fprintf(os.Stderr, "  It uses files from the github/gitignore repository.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults() // Prints descriptions for defined flags (like --list)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage examples:\n")
		fmt.Fprintf(os.Stderr, "  %s <language>        (e.g., %s go, %s python, %s node)\n", os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --list            (To see all available languages)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if doList {
		fileData, err := fetchList(CONTENTS_URL)
		if err != nil {
			fmt.Printf("could not fetch file list form %s\n", CONTENTS_URL)
			os.Exit(1)
		}
		fileList := loadFiles(fileData)
		displayFileList(fileList)
		os.Exit(0)
	}

	language := args[0]
	langUrl, _ := url.JoinPath(
		RAWP_REFIX,
		fmt.Sprintf("%s%s", language, ".gitignore"),
	)

	resp, err := http.Get(langUrl)
	if err != nil {
		fmt.Printf("could not fetch file form %s\n", langUrl)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("received non-OK HTTP status for %s: %s\n", langUrl, resp.Status)
		os.Exit(1)
	}

	out, err := os.Create(".gitignore")
	if err != nil {
		fmt.Printf("could not create .gitignore file:\n%s\n", err)
		os.Exit(1)
	}
	defer out.Close()

	bytesCopied, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("failed to copy content to file %s: %s\n", ".gitignore", err)
	}
	fmt.Printf("Successfully downloaded %d bytes to %s\n", bytesCopied, ".gitignore")

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
	resp, err := http.Get(url)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return content, err
	}
	return content, nil

}

const LIST_HEADER = "Available .gitignore files:\n\n"
const SEP = "---\n"

func displayFileList(files []string) {
	var currentFirst rune
	slices.SortFunc(files, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	fmt.Print(LIST_HEADER)
	for _, f := range files {
		// get first complete rune, not only first byte
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
