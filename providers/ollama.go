package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OllamaProvider implements the Provider interface for Ollama
type OllamaProvider struct {
	URL    string
	Model  string
	Client *http.Client
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
}

// NewOllamaProvider creates a new Ollama provider instance
func NewOllamaProvider(url, model string) *OllamaProvider {
	return &OllamaProvider{
		URL:    url,
		Model:  model,
		Client: &http.Client{},
	}
}

// SendMessage sends a request to Ollama and returns the response
func (p *OllamaProvider) SendMessage(messages []Message, systemPrompt string) (string, error) {
	// Convert to Ollama format and prepend system message
	ollamaMessages := make([]ollamaMessage, 0, len(messages)+1)
	ollamaMessages = append(ollamaMessages, ollamaMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	for _, msg := range messages {
		role := msg.Role
		// Convert "model" role to "assistant" for Ollama
		if role == "model" {
			role = "assistant"
		}
		ollamaMessages = append(ollamaMessages, ollamaMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	reqBody := ollamaChatRequest{
		Model:    p.Model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := p.Client.Post(p.URL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var chatResp ollamaChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return strings.TrimSpace(chatResp.Message.Content), nil
}
