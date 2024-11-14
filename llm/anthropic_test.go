package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liushuangls/go-anthropic"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestClaudeClient_SubmitTask_Success(t *testing.T) {
	// Create a mock server
	mockResponse := `{"completion": "This is a test response from Claude."}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create a new ClaudeClient with the mock server URL
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	client := NewClaudeClient("mock-key", logrus.New(), "claude-v1")
	client.client = anthropic.NewClient("mock-api-key", func(c *anthropic.ClientConfig) {
		c.BaseURL = server.URL
	})

	// Call SubmitTask and validate
	response, err := client.SubmitTask(context.Background(), "Hello, Claude!")
	assert.NoError(t, err)
	assert.Equal(t, "This is a test response from Claude.", response)
}

func TestClaudeClient_SubmitTask_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	// Create a new ClaudeClient with the mock server URL
	client := NewClaudeClient("mock-key", logrus.New(), "claude-v1")
	client.client = anthropic.NewClient("mock-api-key", func(c *anthropic.ClientConfig) {
		c.BaseURL = server.URL
	})

	// Call SubmitTask and expect an error
	response, err := client.SubmitTask(context.Background(), "Hello, Claude!")
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "failed to submit task to Claude")
}
