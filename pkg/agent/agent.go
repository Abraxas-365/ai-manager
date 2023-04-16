package agent

import (
	"strings"

	"github.com/Abraxas-365/ai-manager/pkg/chains"
	"github.com/Abraxas-365/ai-manager/pkg/llm"
	"github.com/Abraxas-365/ai-manager/pkg/prompt"
	"github.com/Abraxas-365/ai-manager/pkg/tools"
)

type Agent struct {
	Llm    llm.LLM
	Chain  chains.Chain
	Tools  []tools.Tool
	Prompt prompt.PromptTemplate
}

func FromLlmAndTools(llm llm.LLM, chain chains.Chain, tools []tools.Tool, prompt prompt.PromptTemplate) Agent {
	return Agent{
		llm,
		chain,
		tools,
		prompt,
	}
}

func GetActionInput(input string) (action, actionInput string) {
	fields := strings.Split(input, "\n")
	for _, field := range fields {
		if strings.HasPrefix(field, "Action: ") {
			action = strings.TrimPrefix(field, "Action: ")
		} else if strings.HasPrefix(field, "Action Input: ") {
			actionInput = strings.TrimPrefix(field, "Action Input: ")
		}
	}
	return action, actionInput
}

func GetObservation(action string, actionInput string, tools []tools.Tool) string {
	for _, tool := range tools {
		if tool.Name() == strings.Trim(action, " ") {
			return "\nObservation: " + tool.Run(actionInput) + "\n"
		}
	}
	return "There is no Tools for the task"
}

func IsAnswer(text string) (bool, string) {
	finalAnswerPrefix := "Final Answer:"
	startIndex := strings.Index(text, finalAnswerPrefix)

	if startIndex == -1 {
		return false, ""
	}

	startIndex += len(finalAnswerPrefix)
	text = text[startIndex:]

	endIndex := strings.Index(text, "\n")

	if endIndex != -1 {
		text = text[:endIndex]
	}

	return true, strings.TrimSpace(text)
}

type AgentInterface interface {
	Run(input string) string
}
