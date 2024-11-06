package llm

import (
	"context"
	"encoding/json"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestSubmitTask_Success tests SubmitTask method for a successful response with multiple choices.
func TestSubmitTask_Success(t *testing.T) {
	// Create a test server that returns a successful response with multiple choices.
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

	// Configure the OpenAI client to use the test server
	clientConfig := openai.ClientConfig{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create GPT instance
	gpt := NewGPT(clientConfig, logger, "gpt-3.5-turbo", "user", []agentv1.Capability{}, []agentv1.Tool{})

	// Call SubmitTask and assert the response
	responses, err := gpt.SubmitTask(context.Background(), "Test prompt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Response 1", "Response 2"}, responses)
}

// TestSubmitTask_NoChoices tests the SubmitTask method when no choices are returned.
func TestSubmitTask_NoChoices(t *testing.T) {
	// Create a test server that returns an empty choices response
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

	clientConfig := openai.ClientConfig{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	gpt := NewGPT(clientConfig, logger, "gpt-3.5-turbo", "user", []agentv1.Capability{}, []agentv1.Tool{})

	// Call SubmitTask and assert an error is returned
	responses, err := gpt.SubmitTask(context.Background(), "Test prompt")
	assert.Error(t, err)
	assert.Nil(t, responses)
	assert.Contains(t, err.Error(), "no response from ChatGPT")
}

// TestSubmitTask_ErrorFromAPI tests SubmitTask when the API returns an error.
func TestSubmitTask_ErrorFromAPI(t *testing.T) {
	// Create a test server that returns a server error
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	clientConfig := openai.ClientConfig{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	gpt := NewGPT(clientConfig, logger, "gpt-3.5-turbo", "user", []agentv1.Capability{}, []agentv1.Tool{})

	// Call SubmitTask and assert an error is returned
	responses, err := gpt.SubmitTask(context.Background(), "Test prompt")
	assert.Error(t, err)
	assert.Nil(t, responses)
	assert.Contains(t, err.Error(), "failed to submit task to ChatGPT")
}

// TestGetCapabilities tests GetCapabilities method.
func TestGetCapabilities(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	capabilities := []agentv1.Capability{
		agentv1.Capability_CAPABILITY_TEXT,
		agentv1.Capability_CAPABILITY_SUMMARIZATION,
	}

	gpt := NewGPT(openai.ClientConfig{}, logger, "gpt-3.5-turbo", "user", capabilities, []agentv1.Tool{})

	// Call GetCapabilities and assert the response
	result := gpt.GetCapabilities()
	assert.Equal(t, capabilities, result)
}

// TestGetTools tests GetTools method.
func TestGetTools(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	tools := []agentv1.Tool{
		{Name: "Tool 1"},
		{Name: "Tool 2"},
	}

	gpt := NewGPT(openai.ClientConfig{}, logger, "gpt-3.5-turbo", "user", []agentv1.Capability{}, tools)

	// Call GetTools and assert the response
	result := gpt.GetTools()
	assert.Equal(t, tools, result)
}

// TestNewGPT tests the NewGPT constructor function.
func TestNewGPT(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	capabilities := []agentv1.Capability{
		agentv1.Capability_CAPABILITY_TEXT,
	}
	tools := []agentv1.Tool{
		{Name: "Tool 1"},
	}

	gpt := NewGPT(openai.ClientConfig{}, logger, "gpt-3.5-turbo", "user", capabilities, tools)

	// Assert that NewGPT initializes correctly
	assert.Equal(t, "gpt-3.5-turbo", gpt.model)
	assert.Equal(t, "user", gpt.role)
	assert.Equal(t, capabilities, gpt.capabilities)
	assert.Equal(t, tools, gpt.tools)
}
