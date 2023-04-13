package openai

import (
	"errors"
	"os"
)

var ErrEmptyResponse = errors.New("no response")
var ErrMissingToken = errors.New("missing the OpenAI API key, set it in the OPENAI_API_KEY environment variable")

type Config struct {
	Temperature float32
	Model       ModelChat
}

type ConfigConstructor struct {
	Temperature float32
	Model       ModelChat
}

func NewConfigConstructor() *ConfigConstructor {
	return &ConfigConstructor{
		Temperature: 0.5,
		Model:       GPT35Turbo,
	}
}

type LLM struct {
	client *Client
	config Config
}

func New(config Config) (*LLM, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingToken
	}
	return &LLM{client: NewClient(apiKey), config: config}, nil
}

func (o *LLM) Generate(prompts []string) ([]string, error) {
	messages := make([]Message, len(prompts))
	results := make([]string, len(prompts))
	for _, prompt := range prompts {
		message := Message{Role: "user", Content: prompt}

		messages = append(messages, message)
		result, err := o.client.GetChatCompletion(messages, o.config.Temperature, o.config.Model)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
		messages = append(messages, Message{Role: "assistant", Content: result})
	}

	return results, nil
}

func (o *LLM) Call(prompt string) (string, error) {
	result, err := o.Generate([]string{prompt})
	if err != nil {
		return "", err
	}
	return result[0], nil
}
