package llm

import (
	"context"
	"fmt"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type GPT struct {
	client       *openai.Client
	logger       *logrus.Entry
	model        string
	role         string
	capabilities []agentv1.Capability
	tools        []agentv1.Tool
}

// SubmitTask sends a task (prompt) to the OpenAI ChatGPT API and returns all responses as a slice of strings.
func (c *GPT) SubmitTask(task string) ([]string, error) {
	logger := c.logger.WithFields(logrus.Fields{
		"task": task})
	logger.Info("Submitting task to ChatGPT")

	// Create a request for the OpenAI API
	req := openai.ChatCompletionRequest{
		Model: c.model, // Use the model specified for this llm
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    c.role, // Use the role specified for this llm
				Content: task,
			},
		},
	}

	// Call the OpenAI API and get the response
	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		logger.WithError(err).Error("Failed to submit task to ChatGPT")
		return nil, fmt.Errorf("failed to submit task to ChatGPT: %w", err)
	}

	// Check if the response contains choices and log them
	if len(resp.Choices) == 0 {
		logger.Warn("No response choices from ChatGPT")
		return nil, fmt.Errorf("no response from ChatGPT")
	}

	var results []string
	for i, choice := range resp.Choices {
		logger.WithFields(logrus.Fields{
			"choice_index": i,
			"content":      choice.Message.Content,
		}).Info("Received choice from ChatGPT")

		results = append(results, choice.Message.Content)
	}

	logger.Info("Task successfully processed by ChatGPT with multiple choices")
	return results, nil
}

// GetCapabilities returns the predefined capabilities of the GPT llm.
func (c *GPT) GetCapabilities() []agentv1.Capability {
	c.logger.Info("Retrieving capabilities of GPT llm")
	return c.capabilities
}

// GetTools returns the predefined tools of the GPT llm.
func (c *GPT) GetTools() []agentv1.Tool {
	c.logger.Info("Retrieving tools of GPT llm")
	return c.tools
}

// NewGPT creates a new instance of GPT with the given OpenAI client configuration, a logger, model, and role.
func NewGPT(config openai.ClientConfig, logger *logrus.Logger, model, role string, capabilities []agentv1.Capability, tools []agentv1.Tool) *GPT {
	client := openai.NewClientWithConfig(config)
	entry := logger.WithFields(logrus.Fields{
		"module": "GPT",
		"model":  model,
		"role":   role},
	)
	return &GPT{
		client:       client,
		logger:       entry,
		model:        model,
		role:         role,
		capabilities: capabilities,
		tools:        tools,
	}
}
