package llm_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/ollama/ollama/api"
	_ "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTask_Success(t *testing.T) {
	// Set up a test HTTP server to mock the Ollama API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure correct request path
		if r.URL.Path != "/api/generate" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		// Simulate a successful response from the API
		response := api.GenerateResponse{
			Response: "Test response",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Configure the Ollama client to use the test server
	logger, _ := test.NewNullLogger()
	parsedURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	client := api.NewClient(parsedURL, http.DefaultClient)
	ollamaClient := llm.NewOllama(client, logger, "test-model", "user", []agentinfo.Capability{}, []agentinfo.Tool{})

	// Call SubmitTask
	o := &agentinfo.LlmOpt{Typ: agentinfo.LlmOptType_LLM_OPT_TYPE_OLLAMA_RESPONSE_FORMAT}
	err = llm.FromVal(o, "json")
	require.NoError(t, err)
	responses, err := ollamaClient.SubmitTask(context.Background(), "Hello", o)
	assert.NoError(t, err)
	assert.Equal(t, "Test response", responses)
}

func TestSubmitTask_Error(t *testing.T) {
	// Set up a test HTTP server to mock an error response from the Ollama API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "API error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Configure the Ollama client to use the test server
	logger, hook := test.NewNullLogger()
	parsedURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	client := api.NewClient(parsedURL, http.DefaultClient)
	ollamaClient := llm.NewOllama(client, logger, "test-model", "user", []agentinfo.Capability{}, []agentinfo.Tool{})

	// Call SubmitTask and expect an error
	o := &agentinfo.LlmOpt{Typ: agentinfo.LlmOptType_LLM_OPT_TYPE_OLLAMA_RESPONSE_FORMAT}
	err = llm.FromVal(o, "json")
	require.NoError(t, err)
	responses, err := ollamaClient.SubmitTask(context.Background(), "Hello", o)
	assert.Error(t, err)
	assert.Empty(t, responses)
	assert.Contains(t, err.Error(), "failed to submit task to Ollama")

	// Verify logs
	assert.Equal(t, "Failed to submit task to Ollama", hook.LastEntry().Message)
}

func TestGetCapabilities(t *testing.T) {
	// Set up logger and ollama client without using the server for this test
	logger, hook := test.NewNullLogger()
	parsedURL, err := url.Parse("http://localhost")
	require.NoError(t, err)
	client := api.NewClient(parsedURL, nil)
	capabilities := []agentinfo.Capability{{Id: "1", Description: "text generation"}}
	ollamaClient := llm.NewOllama(client, logger, "test-model", "user", capabilities, []agentinfo.Tool{})

	// Call GetCapabilities
	retrievedCapabilities := ollamaClient.GetCapabilities()
	assert.Equal(t, capabilities, retrievedCapabilities)

	// Verify logs
	assert.Equal(t, "Retrieving capabilities of Ollama client", hook.LastEntry().Message)
}

func TestGetTools(t *testing.T) {
	// Set up logger and ollama client without using the server for this test
	logger, hook := test.NewNullLogger()
	parsedURL, err := url.Parse("http://localhost")
	require.NoError(t, err)
	client := api.NewClient(parsedURL, nil) // URL isn't relevant here
	tools := []agentinfo.Tool{{Name: "Tool1"}}
	ollamaClient := llm.NewOllama(client, logger, "test-model", "user", []agentinfo.Capability{}, tools)

	// Call GetTools
	retrievedTools := ollamaClient.GetTools()
	assert.Equal(t, tools, retrievedTools)

	// Verify logs
	assert.Equal(t, "Retrieving tools of Ollama client", hook.LastEntry().Message)
}
