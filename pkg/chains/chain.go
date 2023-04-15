package chains

import (
	"github.com/Abraxas-365/ai-manager/pkg/llm"
	"github.com/Abraxas-365/ai-manager/pkg/prompt"
)

type Chain struct {
	llm    llm.LLM
	prompt prompt.PromptTemplate
}

func NewChain(llm llm.LLM, prompt prompt.PromptTemplate) Chain {
	return Chain{llm: llm, prompt: prompt}
}

func (c *Chain) Run(inputVariable string, stop []string) (string, error) {

	completePrompt, err := c.prompt.Format(
		map[string]interface{}{
			c.prompt.InputVariables[0]: inputVariable,
		},
	)
	if err != nil {
		return "", err
	}

	resp, err := c.llm.Call(completePrompt, stop)
	if err != nil {
		return "", err
	}

	return resp, nil

}
