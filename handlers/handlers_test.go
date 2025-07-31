package handlers

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadFiles(t *testing.T) {
	t.Run("TestLoadFiles_0", func(t *testing.T) {
		// Create a mock repo content
		var content RepoContent

		// Call the function with mock client
		result0 := loadFiles(content)

		// Verify the result is empty
		if len(result0) != 0 {
			t.Errorf("Expected empty slice, got %v", result0)
		}
	})
}

func TestFetchList(t *testing.T) {
	t.Run("TestFetchList_0", func(t *testing.T) {
		// Create a mock client with a successful response
		client := NewMockClientWithSuccess(`
			[{"name": "python.gitignore", "type": "file"}]
		`)
		url := ""

		// Call the function with mock client
		result0, result1 := fetchList(client, url)

		// For RepoContent, we need to compare content
		expected := []ContentEntry{
			{Name: "python.gitignore", Type: "file"},
		}
		if len(result0) != len(expected) {
			t.Errorf("Expected %d entries, got %d", len(expected), len(result0))
		}
		for i, entry := range result0 {
			if entry.Name != expected[i].Name || entry.Type != expected[i].Type {
				t.Errorf("Expected entry %v, got %v", expected[i], entry)
			}
		}

		// For errors, we can compare directly
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}

	})
}
func TestDownLoadFile(t *testing.T) {
	t.Run("TestDownLoadFile_0", func(t *testing.T) {
		// Create a mock client with a successful response
		client := NewMockClientWithSuccess("mock file content")
		langUrl := ""

		// Call the function with mock client
		result0, result1 := downLoadFile(client, langUrl)

		// For io.ReadCloser, we can compare directly
		if result0 == nil {
			t.Errorf("Expected non-nil ReadCloser, got nil")
		}

		// For errors, we can compare directly
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}

		// Read the content from the ReadCloser and verify it matches
		content, err := io.ReadAll(result0)
		if err != nil {
			t.Errorf("Expected no error reading from ReadCloser, got %v", err)
		}
		if string(content) != "mock file content" {
			t.Errorf("Expected content 'mock file content', got '%s'", string(content))
		}

	})
}

func TestResourceAvailable(t *testing.T) {
	t.Run("TestResourceAvailable_0", func(t *testing.T) {
		// Create a mock client with a successful response with StatusNotFound
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusNotFound)
		recorder.Body.WriteString("{\"message\": \"Not Found\"}")
		client := &MockClient{
			Response: recorder.Result(),
			Error:    nil,
		}
		url := ""

		// Call the function with mock client
		result0, result1 := resourceAvailable(client, url)

		// For bool, we can compare directly
		if result0 != false {
			t.Errorf("Expected false, got %v", result0)
		}

		// For errors, we can compare directly
		if result1 != ErrGitignoreNotFound {
			t.Errorf("Expected ErrGitignoreNotFound, got %v", result1)
		}

	})
}

func TestPromptForOverwrite(t *testing.T) {
	t.Run("TestPromptForOverwrite_0", func(t *testing.T) {
		// Create mock reader and writer
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		in := bufio.NewReadWriter(bufio.NewReader(&inBuf), bufio.NewWriter(&inBuf))
		out := bufio.NewReadWriter(bufio.NewReader(&outBuf), bufio.NewWriter(&outBuf))

		// Call the function with mock reader and writer
		var result0 choice
		var result1 error
		read := make(chan struct{})
		go func() {
			result0, result1 = promptForOverwrite(in, out)
			read <- struct{}{}
		}()

		// Simulate user input
		inBuf.WriteString("o\n")
		out.Flush()

		// Wait for the function to complete
		<-read

		// Verify the result
		if result0 != ChoiceOverwrite {
			t.Errorf("Expected ChoiceOverwrite, got %v", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}

	})

	t.Run("TestPromptForOverwrite_1", func(t *testing.T) {
		// Create mock reader and writer
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		in := bufio.NewReadWriter(bufio.NewReader(&inBuf), bufio.NewWriter(&inBuf))
		out := bufio.NewReadWriter(bufio.NewReader(&outBuf), bufio.NewWriter(&outBuf))

		// Call the function with mock reader and writer
		var result0 choice
		var result1 error
		read := make(chan struct{})
		go func() {
			result0, result1 = promptForOverwrite(in, out)
			read <- struct{}{}
		}()

		// Simulate user input
		inBuf.WriteString("a\n")
		out.Flush()

		// Wait for the function to complete
		<-read

		// Verify the result
		if result0 != ChoiceAppend {
			t.Errorf("Expected ChoiceAppend, got %v", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}

	})

	t.Run("TestPromptForOverwrite_2", func(t *testing.T) {
		// Create mock reader and writer
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		in := bufio.NewReadWriter(bufio.NewReader(&inBuf), bufio.NewWriter(&inBuf))
		out := bufio.NewReadWriter(bufio.NewReader(&outBuf), bufio.NewWriter(&outBuf))

		// Call the function with mock reader and writer
		var result0 choice
		var result1 error
		read := make(chan struct{})
		go func() {
			result0, result1 = promptForOverwrite(in, out)
			read <- struct{}{}
		}()

		// Simulate user input
		inBuf.WriteString("c\n")
		out.Flush()

		// Wait for the function to complete
		<-read

		// Verify the result
		if result0 != ChoiceCancel {
			t.Errorf("Expected ChoiceCancel, got %v", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}

	})
}
