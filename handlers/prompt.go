package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrOpCancelledByUser = errors.New(
	"Operation cancelled by user",
)

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
