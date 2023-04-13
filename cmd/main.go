package main

import (
	"fmt"
	"strings"

	"github.com/Abraxas-365/ai-manager/pkg/prompt"
)

func main() {
	inputVariables := []string{"test"}
	tmpl := `The available tools are:
{{.tool_names}}

Descriptions:
{{.tool_descriptions}}
	{{.test}}
`

	tools := []struct {
		Name        string
		Description string
	}{
		{"ToolA", "Tool A does something awesome."},
		{"ToolB", "Tool B does something even more awesome."},
	}

	var toolNames []string
	var toolDescriptions []string
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name)
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("%s: %s", tool.Name, tool.Description))
	}

	partialVariables := map[string]interface{}{
		"tool_names":        strings.Join(toolNames, ", "),
		"tool_descriptions": strings.Join(toolDescriptions, "\n"),
	}

	pt := prompt.NewPromptTemplate(inputVariables, tmpl, partialVariables)

	formatted, err := pt.Format(map[string]interface{}{
		"test": "test",
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(formatted)
}
