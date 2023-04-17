package openai

import (
	"errors"
	"net/http"
)

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type Data struct {
	Model       string    `json:"model,omitempty"`
	Messages    []Message `json:"messages,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	Prompt      string    `json:"prompt,omitempty"`
	Stop        []string  `json:"stop,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

var ErrEmptyResponse = errors.New("no response")
var ErrMissingToken = errors.New("missing the OpenAI API key, set it in the OPENAI_API_KEY environment variable")
