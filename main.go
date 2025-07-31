package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/moritz-tiesler/igor/handlers"
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

type Config struct {
	List     bool
	Language string
}

func main() {

	flag.Parse()

	cfg := Config{
		List: doList,
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	run(cfg, flag.Args(), client)

}

func run(cfg Config, args []string, client handlers.Client) {

	if cfg.List {
		err := handlers.List(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching list: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	open := func(path string, flag int, perm os.FileMode) (io.WriteCloser, error) {
		return os.OpenFile(path, flag, perm)
	}

	language := args[0]
	bytesWritten, err := handlers.PullIgnoreFile(
		client,
		language,
		handlers.PromptForOverwrite,
		handlers.Exists,
		open,
	)

	if err != nil {
		if errors.Is(err, handlers.ErrGitignoreNotFound) {
			fmt.Fprintf(os.Stderr, "Error: No .gitignore file found for '%s'.\n", language)
			fmt.Fprintf(os.Stderr, "Please ensure you have typed the language name correctly.\n")
			fmt.Fprintf(os.Stderr, "For a full list of available languages, use: %s --list\n", os.Args[0])
			os.Exit(1)
		}
		if errors.Is(err, handlers.ErrOpCancelledByUser) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error downloading .gitignore file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File downloaded successfully!")
	fmt.Printf("%d bytes written to %s\n", bytesWritten, ".gitignore")
}
