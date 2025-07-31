package handlers

import (
	"io"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("TestList_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var client Client

		result0 := List(client)

		var expect0 error
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

	})
}

func TestLoadFiles(t *testing.T) {
	t.Run("TestLoadFiles_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var content RepoContent

		result0 := loadFiles(content)

		var expect0 []string
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

	})
}

func TestFetchList(t *testing.T) {
	t.Run("TestFetchList_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var client Client
		var url string

		result0, result1 := fetchList(client, url)

		var expect0 RepoContent
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

		var expect1 error
		if !reflect.DeepEqual(result1, expect1) {
			t.Errorf("Expected %v, got %v", expect1, result1)
		}

	})
}

func TestPullIgnoreFile(t *testing.T) {
	t.Run("TestPullIgnoreFile_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var client Client
		var language string

		result0, result1 := PullIgnoreFile(client, language)

		var expect0 int64
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

		var expect1 error
		if !reflect.DeepEqual(result1, expect1) {
			t.Errorf("Expected %v, got %v", expect1, result1)
		}

	})
}

func TestDownLoadFile(t *testing.T) {
	t.Run("TestDownLoadFile_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var client Client
		var langUrl string

		result0, result1 := downLoadFile(client, langUrl)

		var expect0 io.ReadCloser
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

		var expect1 error
		if !reflect.DeepEqual(result1, expect1) {
			t.Errorf("Expected %v, got %v", expect1, result1)
		}

	})
}

func TestResourceAvailable(t *testing.T) {
	t.Run("TestResourceAvailable_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var client Client
		var url string

		result0, result1 := resourceAvailable(client, url)

		var expect0 bool
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

		var expect1 error
		if !reflect.DeepEqual(result1, expect1) {
			t.Errorf("Expected %v, got %v", expect1, result1)
		}

	})
}

func TestPromptForOverwrite(t *testing.T) {
	t.Run("TestPromptForOverwrite_0", func(t *testing.T) {

		// delete this after your implementation
		t.Fatalf("test not implemented")

		var in io.Reader
		var out io.Writer

		result0, result1 := promptForOverwrite(in, out)

		var expect0 choice
		if !reflect.DeepEqual(result0, expect0) {
			t.Errorf("Expected %v, got %v", expect0, result0)
		}

		var expect1 error
		if !reflect.DeepEqual(result1, expect1) {
			t.Errorf("Expected %v, got %v", expect1, result1)
		}

	})
}
