package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type LLMClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	Model      string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	// OpenRouter specific fields
	RouteType  string   `json:"route,omitempty"`
	Transforms []string `json:"transforms,omitempty"`
}

type ChatChoice struct {
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	ID      string       `json:"id"`
	Choices []ChatChoice `json:"choices"`
	Model   string       `json:"model"`
}

// NewOpenRouterClient creates a new LLM client for OpenRouter
func NewOpenRouterClient(model string) *LLMClient {
	if model == "" {
		// Default to a reliable model if none specified
		model = "anthropic/claude-3.7-sonnet"
	}

	return &LLMClient{
		APIKey:     os.Getenv("OPENROUTER_API_KEY"),
		BaseURL:    "https://openrouter.ai/api/v1/chat/completions",
		HTTPClient: &http.Client{},
		Model:      model,
	}
}

// ChatCompletion sends a chat request to the LLM via OpenRouter
func (c *LLMClient) ChatCompletion(messages []ChatMessage) (string, error) {
	reqBody, err := json.Marshal(ChatRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ListAvailableModels returns information about which models are available
func (c *LLMClient) ListAvailableModels() (string, error) {
	// OpenRouter model info endpoint
	modelsURL := "https://openrouter.ai/api/v1/models"

	req, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("HTTP-Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed: %s", body)
	}

	// Just return the raw JSON for now - the app could parse this and display model options
	return string(body), nil
}

