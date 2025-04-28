package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	Stream      bool          `json:"stream,omitempty"`
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

type StreamingChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func NewOpenRouterClient(model string) *LLMClient {
	if model == "" {
		model = "anthropic/claude-3.7-sonnet"
	}

	return &LLMClient{
		APIKey:     os.Getenv("OPENROUTER_API_KEY"),
		BaseURL:    "https://openrouter.ai/api/v1/chat/completions",
		HTTPClient: &http.Client{},
		Model:      model,
	}
}

// ChatCompletionStream streams the response chunks to the provided channel
func (c *LLMClient) ChatCompletionStream(messages []ChatMessage,
	outputChan chan<- string, doneChan chan<- bool) {

	reqBody, err := json.Marshal(ChatRequest{
		Model:       c.Model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
		Stream:      true, // Enable streaming
	})
	if err != nil {
		outputChan <- fmt.Sprintf("Error preparing request: %v", err)
		doneChan <- true
		return
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		outputChan <- fmt.Sprintf("Error creating request: %v", err)
		doneChan <- true
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("X-Title", "OpenAPI Chat Assistant")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		outputChan <- fmt.Sprintf("Error sending request: %v", err)
		doneChan <- true
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		outputChan <- fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, body)
		doneChan <- true
		return
	}

	reader := bufio.NewReaderSize(resp.Body, 32768)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			outputChan <- fmt.Sprintf("\nError reading stream: %v", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}

		// Remove the "data: " prefix that's common in SSE
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		var streamResp StreamingChatResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			continue // Skip malformed chunks
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				outputChan <- content
			}

			// Check if we're done
			if streamResp.Choices[0].FinishReason != "" {
				break
			}
		}
	}

	// Signal that we're done
	doneChan <- true
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
