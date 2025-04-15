package llm

import (
	"context"
	"fmt"
	"github.com/liushuangls/go-anthropic"
	"github.com/sirupsen/logrus"
)

// ClaudeClient is a wrapper around the go-anthropic API client.
type ClaudeClient struct {
	client *anthropic.Client
	logger *logrus.Entry
	model  string
}

// NewClaudeClient initializes a Claude client with API key, logger, model, capabilities, and tools.
func NewClaudeClient(apiKey string, logger *logrus.Logger, model string) *ClaudeClient {
	entry := logger.WithFields(logrus.Fields{
		"module": "ClaudeClient",
		"model":  model,
	})
	return &ClaudeClient{
		client: anthropic.NewClient(apiKey),
		logger: entry,
		model:  model,
	}
}

// SubmitTask sends a task (prompt) to the Claude API and returns the response.
func (c *ClaudeClient) SubmitTask(ctx context.Context, task string, _ any) (string, error) {
	c.logger.WithFields(logrus.Fields{
		"task": task,
	}).Info("Submitting task to Claude")

	// Create a request using the go-anthropic client
	req := anthropic.CompleteRequest{
		Model:  c.model,
		Prompt: task,
	}

	// Call the Claude API and get the response
	resp, err := c.client.CreateComplete(ctx, req)
	if err != nil {
		c.logger.WithError(err).Error("Failed to submit task to Claude")
		return "", fmt.Errorf("failed to submit task to Claude: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"content": resp.Completion,
	}).Info("Received response from Claude")

	return resp.Completion, nil
}
