package agent

type OllamaAgent struct {
	apiKey  string
	baseURL string
}

func (c *OllamaAgent) SubmitTask(task string) (string, error) {
	// Implement the logic for making a request to OpenAI's ChatGPT API
	// Example: Use Go's http.Client to send the request
	panic("implement me")
}

func (c *OllamaAgent) GetCapabilities() ([]string, error) {
	// Return the predefined capabilities of ChatGPT
	return []string{"text-generation", "summarization"}, nil
}

func NewOllamaAgent(apiKey, baseURL string) *OllamaAgent {
	return &OllamaAgent{apiKey: apiKey, baseURL: baseURL}
}
