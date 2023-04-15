package llm

// LLM is a langchaingo Large Language Model.
//TODO create a config struct to send config json variables
type LLM interface {
	Call(prompt string, stop []string) (string, error)
	Generate(prompts []string, stop []string) ([]string, error)
}
