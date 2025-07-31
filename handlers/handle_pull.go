package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	GIT_IGNORE = ".gitignore"
	RAW_PREFIX = "https://raw.githubusercontent.com/github/gitignore/main/"

	APPEND    = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	OVERWRITE = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
}

var ErrGitignoreNotFound = errors.New(
	".gitignore file not found for the specified language",
)

func PullIgnoreFile(
	client Client,
	language string,
	promptForOverwrite func(io.Reader, io.Writer) (choice, error),
	fileCheck func(string) (bool, error),
	openFile func(string, int, os.FileMode) (*os.File, error),
) (int64, error) {
	langUrl, _ := url.JoinPath(
		RAW_PREFIX,
		fmt.Sprintf("%s%s", language, ".gitignore"),
	)

	ok, _ := resourceAvailable(client, langUrl)
	if !ok {
		return 0, ErrGitignoreNotFound
	}

	var shouldAppend bool
	var userAction choice = "overwrite" // Default to overwrite if no file exists
	fExists, err := Exists(GIT_IGNORE)
	if err != nil {
		return 0, fmt.Errorf("failed to check for file %s: %w", GIT_IGNORE, err)
	}
	if fExists {
		// file exists
		choice, err := promptForOverwrite(os.Stdin, os.Stdout)
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
	}

	var fileMode int
	if shouldAppend {
		fileMode = APPEND
	} else {
		fileMode = OVERWRITE
	}

	fmt.Printf("Pulling '%s'...\n", langUrl)
	body, err := downLoadFile(client, langUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to download file %s: %w", langUrl, err)
	}
	defer body.Close()

	out, err := openFile(GIT_IGNORE, fileMode, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open/create file %s: %w", GIT_IGNORE, err)
	}
	defer out.Close()

	var bytesWritten int64
	if fileMode == APPEND {
		n, err := out.WriteString("\n")
		if err != nil {
			return 0, fmt.Errorf("failed to write append separator: %w", err)
		}
		bytesWritten += int64(n)
	}
	n, err := io.Copy(out, body)
	if err != nil {
		return 0, fmt.Errorf("failed to copy content to file %s: %w", GIT_IGNORE, err)
	}
	bytesWritten += n

	return bytesWritten, nil
}

func downLoadFile(client Client, langUrl string) (io.ReadCloser, error) {
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

func resourceAvailable(client Client, url string) (bool, error) {

	resp, err := client.Head(url)
	if err != nil {
		return false, fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, ErrGitignoreNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-OK HTTP status for %s: %s", url, resp.Status)
	}
	return true, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		// some fs or permission error
		return false, fmt.Errorf("failed to check for file %s: %w", path, err)
	}
}
