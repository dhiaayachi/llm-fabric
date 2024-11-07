package llm

import (
	"context"
	"fmt"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
)

// OllamaClient is a wrapper around the Ollama API client.
type OllamaClient struct {
	client       *api.Client
	logger       *logrus.Entry
	model        string
	role         string
	capabilities []agentv1.Capability
	tools        []agentv1.Tool
}

var _ Llm = &OllamaClient{}

// SubmitTask sends a task (prompt) to the Ollama API and returns all responses as a slice of strings.
func (c *OllamaClient) SubmitTask(ctx context.Context, task string, opts ...*Opt) (string, error) {
	logger := c.logger.WithFields(logrus.Fields{
		"task": task,
	})
	logger.Info("Submitting task to Ollama")

	respFormat := getOpt[string](agentv1.LlmOptType_LLM_OPT_TYPE_OllamaResponseFormat, opts...)

	// Create a request using Ollama's client
	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: task,
		Format: respFormat,
	}

	var resp string
	// Call the Ollama API and get the response
	err := c.client.Generate(ctx, req, func(response api.GenerateResponse) error {
		resp += response.Response
		logger.WithField("response", resp).Trace("Got response")
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("Failed to submit task to Ollama")
		return "", fmt.Errorf("failed to submit task to Ollama: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"content": resp,
	}).Info("Received response from Ollama")

	return resp, nil
}

// GetCapabilities returns the predefined capabilities of the Ollama client.
func (c *OllamaClient) GetCapabilities() []agentv1.Capability {
	c.logger.Info("Retrieving capabilities of Ollama client")
	return c.capabilities
}

// GetTools returns the predefined tools of the Ollama client.
func (c *OllamaClient) GetTools() []agentv1.Tool {
	c.logger.Info("Retrieving tools of Ollama client")
	return c.tools
}

// NewOllama creates a new instance of OllamaClient with the given configuration, logger, model, and role.
func NewOllama(apiClient *api.Client, logger *logrus.Logger, model, role string, capabilities []agentv1.Capability, tools []agentv1.Tool) *OllamaClient {
	entry := logger.WithFields(logrus.Fields{
		"module": "OllamaClient",
		"model":  model,
		"role":   role,
	})
	return &OllamaClient{
		client:       apiClient,
		logger:       entry,
		model:        model,
		role:         role,
		capabilities: capabilities,
		tools:        tools,
	}
}
