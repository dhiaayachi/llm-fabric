package llm

import (
	"context"
	"fmt"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// OllamaClient is a wrapper around the langchain Ollama API client.
type OllamaClient struct {
	logger *logrus.Entry
	model  string
	role   string
	url    string
}

func (c *OllamaClient) SubmitTaskWithSchema(ctx context.Context, task string, schema string) (response string, err error) {
	//TODO implement me
	panic("implement me")
}

var _ Llm = &OllamaClient{}

// SubmitTask sends a task (prompt) to the Ollama API and returns all responses as a concatenated string.
func (c *OllamaClient) SubmitTask(ctx context.Context, task string, opts ...*llmoptions.LlmOpt) (string, error) {
	logger := c.logger.WithFields(logrus.Fields{
		"task": task,
	})
	logger.Info("Submitting task to Ollama")

	schema := getOpt[string](llmoptions.LlmOptType_LLM_OPT_TYPE_OLLAMA_RESPONSE_SCHEMA, opts...)

	var llmOpts []ollama.Option
	if c.url != "" {
		urlOpts := ollama.WithServerURL(c.url)
		llmOpts = append(llmOpts, urlOpts)

	}

	if schema != "" {
		llmOpts = append(llmOpts, ollama.WithFormat("json"))
	}

	if c.model != "" {
		llmOpts = append(llmOpts, ollama.WithModel(c.model))
	}

	logger.WithField("schema", schema).WithField("model", c.model).WithField("url", c.url).Info("Submitting task to Ollama With the following options")
	o, err := ollama.New(llmOpts...)
	if err != nil {
		logger.WithError(err).Error("Failed to submit task to Ollama")
		return "", err
	}

	var msgs []llms.MessageContent
	if schema != "" {
		message, err := systemMessage(schema)
		if err != nil {
			logger.WithError(err).Error("Failed to submit task to Ollama")
			return "", fmt.Errorf("failed to submit task to Ollama: %w", err)
		}
		msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeSystem, message))
	}
	msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeHuman, task))

	// Send request using langchain Ollama client
	resp, err := o.GenerateContent(ctx, msgs)
	if err != nil {
		logger.WithError(err).Error("Failed to submit task to Ollama")
		return "", fmt.Errorf("failed to submit task to Ollama: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"content": resp,
	}).Info("Received response from Ollama")

	return resp.Choices[0].Content, nil
}

// NewOllama creates a new instance of OllamaClient with the given configuration, logger, model, and role.
func NewOllama(url string, logger *logrus.Logger, model, role string) *OllamaClient {
	entry := logger.WithFields(logrus.Fields{
		"module": "OllamaClient",
		"model":  model,
		"role":   role,
	})
	return &OllamaClient{
		logger: entry,
		model:  model,
		role:   role,
		url:    url,
	}
}

func systemMessage(schema string) (string, error) {

	return fmt.Sprintf(`Always respond with a JSON object with the following structure: 
%s
`, schema), nil
}
