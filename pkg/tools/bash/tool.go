package bash

import "github.com/Abraxas-365/ai-manager/pkg/tools"

type BashTool struct {
	name        string
	description string
}

func NewBashTool() tools.Tool {
	return &BashTool{
		name: "Bash",
		description: `"Executes commands in a terminal". 
		"Input should be valid Linux command"
		"if you want to write a name use '_' instead of ' '. "
		"Use it  as a helper check if needed to run a tool before this one"
		"Priority: 3"
		`,
	}

}

func (s *BashTool) Name() string {
	return s.name
}

func (s *BashTool) Description() string {
	return s.description
}

func (s *BashTool) Run(query string) string {
	if err := executeBash(query); err != nil {
		return "Comand could not be executed"
	}

	return "Comand executed and task finished with success"
}
