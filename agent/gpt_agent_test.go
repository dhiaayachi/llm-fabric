package agent

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestSubmitTask_Success tests the SubmitTask method with a successful response.
func TestSubmitTask_Success(t *testing.T) {
	// Create a test server that returns a successful response
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{Message: openai.ChatCompletionMessage{Content: "Response 1"}},
				{Message: openai.ChatCompletionMessage{Content: "Response 2"}},
			},
		}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	config := openai.DefaultConfig("")

	config.HTTPClient = server.Client()
	config.BaseURL = server.URL

	agent := NewGPTAgent(config, logger, "gpt-3.5-turbo", "user")

	responses, err := agent.SubmitTask("Test prompt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Response 1", "Response 2"}, responses)
}

// TestSubmitTask_NoChoices tests the SubmitTask method when no choices are returned.
func TestSubmitTask_NoChoices(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{},
		}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := openai.NewClientWithConfig(openai.ClientConfig{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	})

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	agent := &GPTAgent{
		client: client,
		logger: logger.WithField("module", "test"),
		model:  "gpt-3.5-turbo",
		role:   "user",
	}

	responses, err := agent.SubmitTask("Test prompt")
	assert.Error(t, err)
	assert.Nil(t, responses)
}

// TestSubmitTask_ErrorFromAPI tests the SubmitTask method when an API error occurs.
func TestSubmitTask_ErrorFromAPI(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := openai.NewClientWithConfig(openai.ClientConfig{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	})

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	agent := &GPTAgent{
		client: client,
		logger: logger.WithField("module", "test"),
		model:  "gpt-3.5-turbo",
		role:   "user",
	}

	responses, err := agent.SubmitTask("Test prompt")
	assert.Error(t, err)
	assert.Equal(t, "failed to submit task to ChatGPT: error, status code: 500, status: 500 Internal Server Error, message: unexpected end of JSON input, body: ", err.Error())
	assert.Nil(t, responses)
}

// TestGetCapabilities tests the GetCapabilities method.
func TestGetCapabilities(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	agent := &GPTAgent{
		logger: logger.WithField("module", "test"),
		model:  "gpt-3.5-turbo",
		role:   "user",
	}

	capabilities, err := agent.GetCapabilities()
	assert.NoError(t, err)
	assert.Equal(t, []string{"text-generation", "summarization"}, capabilities)
}
