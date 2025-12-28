package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GeminiProvider implements the Provider interface for Google Gemini
type GeminiProvider struct {
	APIKey string
	Model  string
	Client *http.Client
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents          []geminiContent `json:"contents"`
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// NewGeminiProvider creates a new Gemini provider instance
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{
		APIKey: apiKey,
		Model:  model,
		Client: &http.Client{},
	}
}

// SendMessage sends a request to Gemini and returns the response
func (p *GeminiProvider) SendMessage(messages []Message, systemPrompt string) (string, error) {
	if p.APIKey == "" {
		return "", fmt.Errorf("Gemini API key is missing. Please add 'api_key' under 'gemini' section in ~/.config/aish/config.yaml")
	}

	// Convert messages to Gemini format
	contents := make([]geminiContent, 0, len(messages))
	for _, msg := range messages {
		role := msg.Role
		// Convert "assistant" role to "model" for Gemini
		if role == "assistant" {
			role = "model"
		}
		// Skip system messages (handled separately)
		if role == "system" {
			continue
		}
		contents = append(contents, geminiContent{
			Role: role,
			Parts: []geminiPart{
				{Text: msg.Content},
			},
		})
	}

	// Build request with system instruction
	reqBody := geminiRequest{
		Contents: contents,
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{
				{Text: systemPrompt},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		p.Model, p.APIKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Gemini: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text), nil
}
