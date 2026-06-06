package llama

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	llamaHost = flag.String("llama_host", "", "llama.cpp server host:port")
	llamaModel = flag.String("llama_model", "LFM2.5-350M_Q4_K_M.gguf", "Model name for API requests")
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatTemplateKwargs struct {
	EnableThinking bool `json:"enable_thinking"`
}

type chatRequest struct {
	Messages           []chatMessage       `json:"messages"`
	Stream             bool                `json:"stream"`
	ReturnProgress     bool                `json:"return_progress"`
	Model              string              `json:"model"`
	ReasoningFormat    string              `json:"reasoning_format"`
	ChatTemplateKwargs chatTemplateKwargs  `json:"chat_template_kwargs"`
	ReasoningControl   bool                `json:"reasoning_control"`
	BackendSampling    bool                `json:"backend_sampling"`
	TimingsPerToken    bool                `json:"timings_per_token"`
}

type chatResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponseChoice struct {
	FinishReason string              `json:"finish_reason"`
	Index        int                 `json:"index"`
	Message      chatResponseMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatResponseChoice `json:"choices"`
}

func Enabled() bool {
	return *llamaHost != ""
}

func Complete(prompt string) (string, error) {
	if !Enabled() {
		return "", fmt.Errorf("llama: not enabled (--llama_host not set)")
	}

	req := chatRequest{
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		Stream:          false,
		ReturnProgress:  false,
		Model:           *llamaModel,
		ReasoningFormat: "auto",
		ChatTemplateKwargs: chatTemplateKwargs{
			EnableThinking: false,
		},
		ReasoningControl: true,
		BackendSampling:  false,
		TimingsPerToken:  false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("llama: marshal request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post("http://"+*llamaHost+"/v1/chat/completions", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llama: POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("llama: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llama: read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(data, &chatResp); err != nil {
		return "", fmt.Errorf("llama: unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("llama: no choices in response")
	}

	content := chatResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("llama: empty content in response")
	}

	return content, nil
}
