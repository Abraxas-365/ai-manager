package mrkl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Abraxas-365/ai-manager/pkg/agent"
	"github.com/Abraxas-365/ai-manager/pkg/chains"
	"github.com/Abraxas-365/ai-manager/pkg/llm"
	"github.com/Abraxas-365/ai-manager/pkg/prompt"
	"github.com/Abraxas-365/ai-manager/pkg/tools"
)

type ZeroShotAgent struct {
	agent agent.Agent
}

const FinalAnswerAction = "Final Answer:"

func GetActionAndInput(llmOutput string) (string, string) {
	if strings.Contains(llmOutput, FinalAnswerAction) {
		finalAnswer := strings.TrimSpace(strings.Split(llmOutput, FinalAnswerAction)[1])
		return "Final Answer", finalAnswer
	}
	regex := regexp.MustCompile(`Action: (.*?)[\n]*Action Input:[\s]*(.*)`)
	match := regex.FindStringSubmatch(llmOutput)
	if len(match) == 0 {
		panic(fmt.Sprintf("Could not parse LLM output: `%s`", llmOutput))
	}
	action := strings.TrimSpace(match[1])
	actionInput := strings.TrimSpace(match[2])
	return action, strings.Trim(actionInput, "\"")
}

func createPrompt(
	tools []tools.Tool,
) prompt.PromptTemplate {
	var toolStrings []string
	for _, tool := range tools {
		toolStrings = append(toolStrings, fmt.Sprintf("%s: %s", tool.Name(), tool.Description()))
	}

	toolString := strings.Join(toolStrings, "\n")
	var toolNames []string
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name())
	}
	toolNamesString := strings.Join(toolNames, ", ")
	formatInstructions := strings.Replace(FORMAT_INSTRUCTIONS, "{tool_names}", toolNamesString, -1)

	template := strings.Join([]string{PREFIX, toolString, formatInstructions, SUFFIX}, "\n\n")

	return prompt.NewPromptTemplateBuilder().AddTemplate(template).AddInputVariables([]string{"input", "agent_scratchpad"}).Build()

}

func FromLlmAndTools(
	llm llm.LLM,
	tools []tools.Tool,
) agent.AgentInterface {

	promptTemplate := createPrompt(tools)
	chain := chains.NewChain(llm, promptTemplate)
	agent := agent.FromLlmAndTools(llm, chain, tools, promptTemplate)
	return &ZeroShotAgent{
		agent,
	}
}

func (z *ZeroShotAgent) Run(input string) string {

	output, _ := z.agent.Chain.Run(input, []string{"\nObservation:", "\n\tObservation"}) //esto va a devolver el lo que continua del tough + action
	for true {
		if isAnswer, answer := agent.IsAnswer(output); isAnswer {
			return answer
		}
		//pasos intermedios
		action, actionInput := agent.GetActionInput(output)
		observation := agent.GetObservation(action, actionInput, z.agent.Tools)
		scratchpad := output + observation
		fmt.Println("step", output)

		//crear el nuevo promp que contenga el scratchpad
		newPrompt := prompt.NewPromptTemplateBuilder().AddTemplate(z.agent.Prompt.Template).AddPartialVariables(map[string]interface{}{
			"input": input,
		}).AddInputVariables([]string{"agent_scratchpad"}).Build()
		//creat un chain y enviarlo
		chain := chains.NewChain(z.agent.Llm, newPrompt)
		output, _ = chain.Run(scratchpad, []string{"\nObservation:", "\n\tObservation"})
	}
	return ""
}
