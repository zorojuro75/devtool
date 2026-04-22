package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Completer is the interface all LLM providers must implement.
// Using an interface means tests can inject a mock — no real HTTP needed.
type Completer interface {
	Stream(ctx context.Context, prompt, system string) (io.Reader, error)
}

type openRouterClient struct {
	apiKey  string
	model   string
	timeout time.Duration
}

func NewOpenRouter(apiKey, model string, timeoutSecs int) Completer {
	return &openRouterClient{
		apiKey:  apiKey,
		model:   model,
		timeout: time.Duration(timeoutSecs) * time.Second,
	}
}

type orRequest struct {
	Model    string     `json:"model"`
	Messages []orMsg    `json:"messages"`
	Stream   bool       `json:"stream"`
}

type orMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *openRouterClient) Stream(ctx context.Context, prompt, system string) (io.Reader, error) {
	body, _ := json.Marshal(orRequest{
		Model: c.model,
		Messages: []orMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: prompt},
		},
		Stream: true,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/zorojuro75/devtool")

	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(b))
	}

	return resp.Body, nil
}