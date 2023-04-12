package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

type PromptTemplate struct {
	Template       string
	InputVariables []string
}

func NewPromptTemplate(template string, inputVariables []string) *PromptTemplate {
	return &PromptTemplate{
		Template:       template,
		InputVariables: inputVariables,
	}
}

func (pt *PromptTemplate) FormatTemplate() (string, error) {
	// Create a regex pattern to match the placeholder format
	placeholderPattern := regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

	for _, variable := range pt.InputVariables {
		// Find the next placeholder in the template
		match := placeholderPattern.FindStringSubmatch(pt.Template)
		if match == nil {
			return "", fmt.Errorf("insufficient placeholders in template")
		}
		placeholder := match[0]

		// Replace the placeholder with the variable in the template
		pt.Template = strings.Replace(pt.Template, placeholder, variable, 1)
	}

	// Check if there are any remaining placeholders in the template
	if placeholderPattern.MatchString(pt.Template) {
		return "", fmt.Errorf("insufficient input variables for template")
	}

	return pt.Template, nil
}
