package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OllamaClient handles communication with Ollama
type OllamaClient struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// OllamaRequest represents the request to Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient() *OllamaClient {
	baseURL := os.Getenv("OLLAMA_URL")
	if baseURL == "" {
		baseURL = "http://192.168.5.10:11434"
	}
	
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "qwen2.5-coder:3b"
	}
	
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 120 * time.Second, // AI responses can take time
		},
	}
}

// Generate sends a prompt to Ollama and returns the response
func (c *OllamaClient) Generate(prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	
	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	return ollamaResp.Response, nil
}

// HealthCheck checks if Ollama is accessible
func (c *OllamaClient) HealthCheck() error {
	url := fmt.Sprintf("%s/api/tags", c.BaseURL)
	resp, err := c.Client.Get(url)
	if err != nil {
		return fmt.Errorf("ollama is not accessible: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// Global Ollama client instance
var Client *OllamaClient

// Initialize sets up the global Ollama client
func Initialize() {
	Client = NewOllamaClient()
}
