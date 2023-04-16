package main

import (
	"fmt"

	"github.com/Abraxas-365/ai-manager/pkg/agent/mrkl"
	"github.com/Abraxas-365/ai-manager/pkg/openai"
	"github.com/Abraxas-365/ai-manager/pkg/tools"
	"github.com/Abraxas-365/ai-manager/pkg/tools/googlesearch"
)

func main() {
	//Declaro que Large lenguage model quiero usar
	llm, err := openai.NewCompletition(
		openai.NewCompletitonConfigConstructor().
			AddMaxTokens(500).
			AddModel(openai.TextDavinchi3).
			AddTemperature(0).
			Build(),
	)
	if err != nil {
		fmt.Println(err)
	}

	//Declaro las herramientas que voy a usar,
	//Wrapper de api de google search
	googleSearchTool, err := googlesearch.NewSearchTool()
	if err != nil {
		fmt.Println(err)
	}
	//Declaro el agente con las herramientas que tiene a su disposicion
	//No necesariemanete las va a usar
	agent := mrkl.FromLlmAndTools(llm, []tools.Tool{googleSearchTool})
	if err != nil {
		fmt.Println(err)
	}
	//Hacerle la pregunta al agente
	//Pregunta de actualidad, peru a cambiado de presidente muchas veces
	//Si le preguntamos a ChatGPT nos va a decir Pedro castillo
	answer := agent.Run("who is the president of the peru and when was she born?")
	fmt.Println("\n",
		"ESta es la respuesta de la IA: ",
		answer)
}
