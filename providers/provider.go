package providers

// Message represents a chat message with role and content
type Message struct {
	Role    string
	Content string
}

// Provider defines the interface for LLM providers
type Provider interface {
	// SendMessage sends messages to the LLM and returns the response
	SendMessage(messages []Message, systemPrompt string) (string, error)
}
