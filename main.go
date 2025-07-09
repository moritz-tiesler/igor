package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
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
	RAW_REFIX    = "https://raw.githubusercontent.com/github/gitignore/main/"
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

	if doList {
		fileData, err := fetchList(CONTENTS_URL)
		if err != nil {
			fmt.Printf("could not fetch file list from %s\n", CONTENTS_URL)
			os.Exit(1)
		}
		fileList := loadFiles(fileData)
		displayFileList(fileList)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	language := args[0]
	langUrl, _ := url.JoinPath(
		RAW_REFIX,
		fmt.Sprintf("%s%s", language, ".gitignore"),
	)

	outputFileName := ".gitignore"

	var shouldAppend bool
	var userAction string = "overwrite" // Default to overwrite if no file exists
	_, err := os.Stat(outputFileName)
	if err == nil {
		choice, promptErr := promptForOverwrite(outputFileName)
		if promptErr != nil {
			fmt.Fprintf(os.Stderr, "Error reading user input: %v\n", promptErr)
			os.Exit(1)
		}
		userAction = choice
		switch userAction {
		case "cancel":
			fmt.Println("Operation cancelled by user.")
			os.Exit(0)
		case "append":
			shouldAppend = true
		case "overwrite":
			fmt.Printf("Overwriting '%s'...\n", outputFileName)
			shouldAppend = false
		}

	} else if !errors.Is(err, fs.ErrNotExist) { // Some error other than "file not found"
		fmt.Fprintf(os.Stderr, "Error checking file status: %v\n", err)
		os.Exit(1)
	}
	err = downLoadFile(langUrl, outputFileName, shouldAppend)
	if err != nil {
		if errors.Is(err, ErrGitignoreNotFound) {
			fmt.Fprintf(os.Stderr, "Error: No .gitignore file found for '%s'.\n", language)
			fmt.Fprintf(os.Stderr, "Please ensure you have typed the language name correctly.\n")
			fmt.Fprintf(os.Stderr, "For a full list of available languages, use: %s --list\n", os.Args[0])
		} else {
			fmt.Fprintf(os.Stderr, "Error downloading .gitignore file: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("File downloaded successfully!")
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

var ErrGitignoreNotFound = errors.New("gitignore file not found for the specified language")

func downLoadFile(langUrl string, filePath string, appendMode bool) error {

	resp, err := http.Get(langUrl)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request to %s: %w", langUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrGitignoreNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status for %s: %s", langUrl, resp.Status)
	}
	// Determine file opening mode
	var fileMode int
	if appendMode {
		fileMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		fileMode = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}

	out, err := os.OpenFile(filePath, fileMode, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create file %s: %w", filePath, err)
	}
	defer out.Close()
	if appendMode {
		if _, err := out.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write append separator: %w", err)
		}
	}

	bytesCopied, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy content to file %s: %w", filePath, err)
	}
	fmt.Printf("%d bytes written to %s\n", bytesCopied, filePath)
	return nil
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

func promptForOverwrite(filePath string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("A '%s' file already exists. What would you like to do? (o)verwrite / (a)ppend / (c)ancel: ", filePath)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input)) // Clean and normalize input

		switch input {
		case "o", "overwrite":
			return "overwrite", nil
		case "a", "append":
			return "append", nil
		case "c", "cancel":
			return "cancel", nil
		default:
			fmt.Println("Invalid choice. Please enter 'o', 'a', or 'c'.")
		}
	}
}
