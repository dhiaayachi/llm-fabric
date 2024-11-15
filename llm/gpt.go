package llm

import (
	"context"
	"fmt"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/sirupsen/logrus"
)

type GPT struct {
	client *openai.Client
	logger *logrus.Entry
	model  string
	role   string
}

var _ Llm = &GPT{}

// SubmitTask sends a task (prompt) to the OpenAI ChatGPT API and returns all responses as a slice of strings.
func (c *GPT) SubmitTask(ctx context.Context, task string, opts ...*llmoptions.LlmOpt) (string, error) {
	logger := c.logger.WithFields(logrus.Fields{
		"task": task})
	logger.Info("Submitting task to ChatGPT")

	schema := getOpt[*jsonschema.Definition](llmoptions.LlmOptType_LLM_OPT_TYPE_GPT_RESPONSE_FORMAT, opts...)
	var rspFormat *openai.ChatCompletionResponseFormat
	if schema != nil {
		rspFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "capabilities",
				Schema: schema,
				Strict: true,
			}}

	}
	// Create a request for the OpenAI API
	req := openai.ChatCompletionRequest{
		Model: c.model, // Use the model specified for this llm
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    c.role, // Use the role specified for this llm
				Content: task,
			},
		},
		ResponseFormat: rspFormat,
	}

	// Call the OpenAI API and get the response
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		logger.WithError(err).Error("Failed to submit task to ChatGPT")
		return "", fmt.Errorf("failed to submit task to ChatGPT: %w", err)
	}

	// Check if the response contains choices and log them
	if len(resp.Choices) == 0 {
		logger.Warn("No response choices from ChatGPT")
		return "", fmt.Errorf("no response from ChatGPT")
	}

	var results string
	for i, choice := range resp.Choices {
		logger.WithFields(logrus.Fields{
			"choice_index": i,
			"content":      choice.Message.Content,
		}).Info("Received choice from ChatGPT")

		results += choice.Message.Content
	}

	logger.Info("Task successfully processed by ChatGPT with multiple choices")
	return results, nil
}

// NewGPT creates a new instance of GPT with the given OpenAI client configuration, a logger, model, and role.
func NewGPT(config openai.ClientConfig, logger *logrus.Logger, model, role string) *GPT {
	client := openai.NewClientWithConfig(config)
	entry := logger.WithFields(logrus.Fields{
		"module": "GPT",
		"model":  model,
		"role":   role},
	)
	return &GPT{
		client: client,
		logger: entry,
		model:  model,
		role:   role,
	}
}
