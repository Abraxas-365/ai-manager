package bash

import "github.com/Abraxas-365/ai-manager/pkg/tools"

type BashTool struct {
	name        string
	description string
}

func NewBashTool() tools.Tool {
	return &BashTool{
		name: "Bash",
		description: `"A guide for common Bash commands and their usage. "
		"Useful for when you need assistance with navigating or manipulating files and directories in a Unix-based system."
		"Input should be a command or a description of the task you want to perform."
		"If you want to create a space-separated forder use '_' instead of spaces"
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
