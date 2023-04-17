package openai

import (
	"os"
)

type CompletitionConfig struct {
	temperature float32
	model       ModelCompletition
}

type CompletitionConfigConstructor struct {
	temperature float32
	model       ModelCompletition
	maxTokens   int
}
type completitionConfigConstructor interface {
	AddTemperature(float32) completitionConfigConstructor
	AddModel(ModelCompletition) completitionConfigConstructor
	Build() CompletitionConfig
}

func NewCompletitonConfigConstructor() *CompletitionConfigConstructor {
	return &CompletitionConfigConstructor{
		temperature: 0,
		model:       TextDavinchi3,
		maxTokens:   1000,
	}
}

func (c *CompletitionConfigConstructor) AddTemperature(temperature float32) completitionConfigConstructor {
	c.temperature = temperature
	return c
}
func (c *CompletitionConfigConstructor) AddMaxTokens(maxTokens int) completitionConfigConstructor {
	c.maxTokens = maxTokens
	return c
}
func (c *CompletitionConfigConstructor) AddModel(model ModelCompletition) completitionConfigConstructor {
	c.model = model
	return c
}
func (c *CompletitionConfigConstructor) Build() CompletitionConfig {
	return CompletitionConfig{
		temperature: c.temperature,
		model:       c.model,
	}
}

type LLMCompletition struct {
	client *Client
	config CompletitionConfig
}

func NewCompletition(config CompletitionConfig) (*LLMCompletition, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingToken
	}
	return &LLMCompletition{client: NewClient(apiKey), config: config}, nil
}

func (o *LLMCompletition) Generate(prompts []string, stop []string) ([]string, error) {
	var completitions []string
	for _, prompt := range prompts {
		result, err := o.client.getCompletion(prompt, 60, o.config.temperature, o.config.model, stop)
		if err != nil {
			return nil, err
		}
		completitions = append(completitions, result)
	}

	return completitions, nil
}

func (o *LLMCompletition) Call(prompt string, stop []string) (string, error) {
	result, err := o.Generate([]string{prompt}, stop)
	if err != nil {
		return "", err
	}
	return result[0], nil
}
