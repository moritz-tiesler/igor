package handlers

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestPullIgnoreFile(t *testing.T) {
	t.Run("TestPullIgnoreFile_0", func(t *testing.T) {
		// Mock dependencies
		var calledPrompt bool
		var calledStat bool
		var calledOpenFile bool

		// Mock client with successful response
		client := NewMockClientWithSuccess("mock content")

		// Create a mock file info for os.Stat

		// Define test function

		// Test case: file exists, user chooses overwrite
		stat := func(string) (bool, error) {
			return true, nil
		}
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		// in := bufio.NewReader(&inBuf)
		out := bufio.NewWriter(&outBuf)

		// Mock promptForOverwrite to return ChoiceOverwrite
		mockPrompt := func(in io.Reader, out io.Writer) (choice, error) {
			calledPrompt = true
			return ChoiceOverwrite, nil
		}

		// Mock openFile to return a mock file
		mockOpenFile := func(path string, mode int, perm os.FileMode) (*os.File, error) {
			calledOpenFile = true
			return &os.File{}, nil
		}

		// Simulate user input
		inBuf.WriteString("o\n")
		out.Flush()

		// Call the function
		result0, result1 := PullIgnoreFile(client, "python", mockPrompt, stat, mockOpenFile)

		// Verify all mocks were called
		if !calledPrompt {
			t.Error("promptForOverwrite was not called")
		}
		if !calledStat {
			t.Error("os.Stat was not called")
		}
		if !calledOpenFile {
			t.Error("os.OpenFile was not called")
		}

		// Verify result
		if result0 != 0 {
			t.Errorf("Expected 0 bytes written, got %d", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}
	})

	t.Run("TestPullIgnoreFile_1", func(t *testing.T) {
		// Mock dependencies
		var calledPrompt bool
		var calledStat bool
		var calledOpenFile bool

		// Mock client with successful response
		client := NewMockClientWithSuccess("mock content")

		// Create a mock file info for os.Stat
		stat := func(string) (bool, error) {
			return true, nil
		}

		// Define test function

		// Test case: file exists, user chooses append
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		// in := bufio.NewReader(&inBuf)
		out := bufio.NewWriter(&outBuf)

		// Mock os.Stat to return nil (file exists)

		// Mock promptForOverwrite to return ChoiceAppend
		mockPrompt := func(in io.Reader, out io.Writer) (choice, error) {
			calledPrompt = true
			return ChoiceAppend, nil
		}

		// Mock openFile to return a mock file
		mockOpenFile := func(path string, mode int, perm os.FileMode) (*os.File, error) {
			calledOpenFile = true
			return &os.File{}, nil
		}

		// Simulate user input
		inBuf.WriteString("a\n")
		out.Flush()

		// Call the function
		result0, result1 := PullIgnoreFile(client, "python", mockPrompt, stat, mockOpenFile)

		// Verify all mocks were called
		if !calledPrompt {
			t.Error("promptForOverwrite was not called")
		}
		if !calledStat {
			t.Error("os.Stat was not called")
		}
		if !calledOpenFile {
			t.Error("os.OpenFile was not called")
		}

		// Verify result
		if result0 != 0 {
			t.Errorf("Expected 0 bytes written, got %d", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}
	})

	t.Run("TestPullIgnoreFile_2", func(t *testing.T) {
		// Mock dependencies
		var calledPrompt bool
		var calledStat bool
		var calledOpenFile bool

		// Mock client with successful response
		client := NewMockClientWithSuccess("mock content")

		// Test case: file does not exist
		stat := func(string) (bool, error) {
			return false, nil
		}
		var inBuf bytes.Buffer
		var outBuf bytes.Buffer
		// in := bufio.NewReader(&inBuf)
		out := bufio.NewWriter(&outBuf)

		// Mock promptForOverwrite to return ChoiceOverwrite
		mockPrompt := func(in io.Reader, out io.Writer) (choice, error) {
			calledPrompt = true
			return ChoiceOverwrite, nil
		}

		// Mock openFile to return a mock file
		mockOpenFile := func(path string, mode int, perm os.FileMode) (*os.File, error) {
			calledOpenFile = true
			return &os.File{}, nil
		}

		// Simulate user input
		inBuf.WriteString("o\n")
		out.Flush()

		// Call the function
		result0, result1 := PullIgnoreFile(client, "python", mockPrompt, stat, mockOpenFile)

		// Verify all mocks were called
		if !calledPrompt {
			t.Error("promptForOverwrite was not called")
		}
		if !calledStat {
			t.Error("os.Stat was not called")
		}
		if !calledOpenFile {
			t.Error("os.OpenFile was not called")
		}

		// Verify result
		if result0 != 0 {
			t.Errorf("Expected 0 bytes written, got %d", result0)
		}
		if result1 != nil {
			t.Errorf("Expected nil error, got %v", result1)
		}
	})
}
