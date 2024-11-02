package agent

type ClaudeAgent struct {
	apiKey  string
	baseURL string
}

func (c *ClaudeAgent) SubmitTask(task string) (string, error) {
	// Implement the logic for making a request to OpenAI's ChatGPT API
	// Example: Use Go's http.Client to send the request
	panic("implement me")
}

func (c *ClaudeAgent) GetCapabilities() ([]string, error) {
	// Return the predefined capabilities of ChatGPT
	return []string{"text-generation", "summarization"}, nil
}

func NewClaudeAgent(apiKey, baseURL string) *ClaudeAgent {
	return &ClaudeAgent{apiKey: apiKey, baseURL: baseURL}
}
