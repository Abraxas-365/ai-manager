package openai

import (
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatConfig struct {
	temperature float32
	model       ModelChat
}

type ChatConfigConstructor struct {
	temperature float32
	model       ModelChat
}
type chatConfigConstructor interface {
	AddTemperature(float32) chatConfigConstructor
	AddModel(ModelChat) chatConfigConstructor
	Build() ChatConfig
}

func NewChatConfigConstructor() *ChatConfigConstructor {
	return &ChatConfigConstructor{
		temperature: 0,
		model:       GPT35Turbo,
	}
}

func (c *ChatConfigConstructor) AddTemperature(temperature float32) chatConfigConstructor {
	c.temperature = temperature
	return c
}
func (c *ChatConfigConstructor) AddModel(model ModelChat) chatConfigConstructor {
	c.model = model
	return c
}
func (c *ChatConfigConstructor) Build() ChatConfig {
	return ChatConfig{
		temperature: c.temperature,
		model:       c.model,
	}
}

type LLMChat struct {
	client *Client
	config ChatConfig
}

func NewChat(config ChatConfig) (*LLMChat, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingToken
	}
	return &LLMChat{client: NewClient(apiKey), config: config}, nil
}

func (o *LLMChat) Generate(prompts []string, stop []string) ([]string, error) {
	messages := []Message{}
	completitions := []string{}
	for _, prompt := range prompts {
		message := Message{Role: "user", Content: prompt}
		messages = append(messages, message)
		result, err := o.client.getChatCompletion(messages, o.config.temperature, o.config.model, stop)
		if err != nil {
			return nil, err
		}
		messages = append(messages, Message{Role: "assistant", Content: result})
		completitions = append(completitions, result)
	}

	return completitions, nil
}

func (o *LLMChat) Call(prompt string, stop []string) (string, error) {
	result, err := o.Generate([]string{prompt}, stop)
	if err != nil {
		return "", err
	}
	return result[0], nil
}
