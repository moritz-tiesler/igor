package handlers

import (
	"net/http"
	"net/http/httptest"
)

// MockClient is a mock implementation of the Client interface for testing
type MockClient struct {
	// Response to return from Get
	Response *http.Response
	// Error to return from Get
	Error error
}

// Get implements the Client interface
func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

// Head implements the Client interface
func (m *MockClient) Head(url string) (*http.Response, error) {
	return m.Response, m.Error
}

// NewMockClient creates a new mock client with the given response and error
func NewMockClient(response *http.Response, err error) *MockClient {
	return &MockClient{
		Response: response,
		Error:    err,
	}
}

// NewMockClientWithSuccess creates a mock client that returns a successful response
func NewMockClientWithSuccess(body string) *MockClient {
	recorder := httptest.NewRecorder()
	recorder.Body.WriteString(body)
	return &MockClient{
		Response: recorder.Result(),
		Error:    nil,
	}
}

// NewMockClientWithError creates a mock client that returns an error
func NewMockClientWithError(err error) *MockClient {
	return &MockClient{
		Response: nil,
		Error:    err,
	}
}
