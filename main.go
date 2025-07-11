package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// flags
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

const (
	CONTENTS_URL = "https://api.github.com/repos/github/gitignore/contents/"
	RAW_PREFIX   = "https://raw.githubusercontent.com/github/gitignore/main/"

	APPEND    = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	OVERWRITE = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	TYPE_FILE = "file"
	TYPE_DIR  = "dir"

	GIT_IGNORE = ".gitignore"

	LIST_HEADER = "Available .gitignore files:\n\n"
	SEP         = "---\n"
)

var ErrGitignoreNotFound = errors.New(
	".gitignore file not found for the specified language",
)

var ErrOpCancelledByUser = errors.New(
	"Operation cancelled by user",
)

func main() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	flag.Parse()

	if doList {
		err := handleDoList(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetich list: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	language := args[0]
	bytesWritten, err := handleFilePull(client, language)

	if err != nil {
		if errors.Is(err, ErrGitignoreNotFound) {
			fmt.Fprintf(os.Stderr, "Error: No .gitignore file found for '%s'.\n", language)
			fmt.Fprintf(os.Stderr, "Please ensure you have typed the language name correctly.\n")
			fmt.Fprintf(os.Stderr, "For a full list of available languages, use: %s --list\n", os.Args[0])
			os.Exit(1)
		}
		if errors.Is(err, ErrOpCancelledByUser) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error downloading .gitignore file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File downloaded successfully!")
	fmt.Printf("%d bytes written to %s\n", bytesWritten, GIT_IGNORE)
}

func handleDoList(client *http.Client) error {
	fileData, err := fetchList(client, CONTENTS_URL)
	if err != nil {
		return fmt.Errorf("could not fetch file list from %s: %w", CONTENTS_URL, err)
	}
	fileList := loadFiles(fileData)
	displayFileList(fileList)
	return nil
}

func handleFilePull(client *http.Client, language string) (int64, error) {
	langUrl, _ := url.JoinPath(
		RAW_PREFIX,
		fmt.Sprintf("%s%s", language, ".gitignore"),
	)

	// todo: check if available first, then prompt, only the do download
	fmt.Printf("Pulling '%s'...\n", langUrl)
	body, err := downLoadFile(client, langUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to download file %s: %w", langUrl, err)
	}
	defer body.Close()

	var shouldAppend bool
	var userAction choice = "overwrite" // Default to overwrite if no file exists
	_, err = os.Stat(GIT_IGNORE)
	if err == nil {
		// file exists
		choice, err := promptForOverwrite()
		if err != nil {
			return 0, fmt.Errorf("Error reading user input: %w\n", err)
		}
		userAction = choice
		switch userAction {
		case ChoiceCancel:
			return 0, ErrOpCancelledByUser
		case ChoiceAppend:
			fmt.Printf("Appending to '%s'...\n", GIT_IGNORE)
			shouldAppend = true
		case ChoiceOverwrite:
			fmt.Printf("Overwriting '%s'...\n", GIT_IGNORE)
			shouldAppend = false
		}
	} else {
		if os.IsNotExist(err) {
			// file will be created
		} else {
			// some fs or permission error
			return 0, fmt.Errorf("failed to check for file %s: %w", GIT_IGNORE, err)
		}
	}

	var fileMode int
	if shouldAppend {
		fileMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		fileMode = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}

	out, err := os.OpenFile(GIT_IGNORE, fileMode, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open/create file %s: %w", GIT_IGNORE, err)
	}
	defer out.Close()

	bytesWritten, err := writeGitIgnore(out, body)
	if fileMode == APPEND {
		if _, err := out.WriteString("\n"); err != nil {
			return 0, fmt.Errorf("failed to write append separator: %w", err)
		}
	}
	if err != nil {
		return bytesWritten, fmt.Errorf("failed to write file: %w", err)
	}

	return bytesWritten, nil
}

func downLoadFile(client *http.Client, langUrl string) (io.ReadCloser, error) {
	resp, err := client.Get(langUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to %s: %w", langUrl, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrGitignoreNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status for %s: %s", langUrl, resp.Status)
	}

	return resp.Body, nil
}

func writeGitIgnore(dest io.Writer, src io.Reader) (int64, error) {
	bytesCopied, err := io.Copy(dest, src)
	if err != nil {
		return 0, fmt.Errorf("failed to copy content to file %s: %w", GIT_IGNORE, err)
	}

	return bytesCopied, nil
}

func loadFiles(content Content) []string {
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

type Content []ContentEntry

type ContentEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func fetchList(client *http.Client, url string) (Content, error) {
	var content Content
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

type choice string

const (
	ChoiceOverwrite choice = "overwrite"
	ChoiceAppend    choice = "append"
	ChoiceCancel    choice = "cancel "
)

func promptForOverwrite() (choice, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("A '%s' file already exists. What would you like to do? (o)verwrite / (a)ppend / (c)ancel: ", GIT_IGNORE)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read user input: %w", err)
		}
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "o", "overwrite":
			return ChoiceOverwrite, nil
		case "a", "append":
			return ChoiceAppend, nil
		case "c", "cancel":
			return ChoiceCancel, nil
		default:
			fmt.Println("Invalid choice. Please enter 'o', 'a', or 'c'.")
		}
	}
}
