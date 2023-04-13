package llm

// LLM is a langchaingo Large Language Model.
type LLM interface {
	Call(prompt string) (string, error)
	Generate(prompts []string) ([]string, error)
}
