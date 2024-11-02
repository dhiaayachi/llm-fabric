package agent

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type GPTAgent struct {
	client *openai.Client
	logger *logrus.Entry
	model  string
	role   string
}

// SubmitTask sends a task (prompt) to the OpenAI ChatGPT API and returns all responses as a slice of strings.
func (c *GPTAgent) SubmitTask(task string) ([]string, error) {
	c.logger.WithFields(logrus.Fields{
		"task":  task,
		"model": c.model,
		"role":  c.role,
	}).Info("Submitting task to ChatGPT")

	// Create a request for the OpenAI API
	req := openai.ChatCompletionRequest{
		Model: c.model, // Use the model specified for this agent
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    c.role, // Use the role specified for this agent
				Content: task,
			},
		},
	}

	// Call the OpenAI API and get the response
	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.logger.WithError(err).Error("Failed to submit task to ChatGPT")
		return nil, fmt.Errorf("failed to submit task to ChatGPT: %w", err)
	}

	// Check if the response contains choices and log them
	if len(resp.Choices) == 0 {
		c.logger.Warn("No response choices from ChatGPT")
		return nil, fmt.Errorf("no response from ChatGPT")
	}

	var results []string
	for i, choice := range resp.Choices {
		c.logger.WithFields(logrus.Fields{
			"choice_index": i,
			"content":      choice.Message.Content,
		}).Info("Received choice from ChatGPT")

		results = append(results, choice.Message.Content)
	}

	c.logger.Info("Task successfully processed by ChatGPT with multiple choices")
	return results, nil
}

// GetCapabilities returns the predefined capabilities of the GPT agent.
func (c *GPTAgent) GetCapabilities() ([]string, error) {
	c.logger.WithFields(logrus.Fields{
		"model": c.model,
		"role":  c.role,
	}).Info("Retrieving capabilities of GPT agent")
	return []string{"text-generation", "summarization"}, nil
}

// NewGPTAgent creates a new instance of GPTAgent with the given OpenAI client configuration, a logger, model, and role.
func NewGPTAgent(config openai.ClientConfig, logger *logrus.Logger, model, role string) *GPTAgent {
	client := openai.NewClientWithConfig(config)
	entry := logger.WithField("module", "GPTAgent")
	return &GPTAgent{
		client: client,
		logger: entry,
		model:  model,
		role:   role,
	}
}
