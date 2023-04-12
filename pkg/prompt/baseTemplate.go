package prompt

import (
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

type PromptTemplate struct {
	InputVariables   []string
	Template         string
	TemplateFormat   string
	ValidateTemplate bool
	PartialVariables []string
}

func (p *PromptTemplate) Format(kwargs map[string]interface{}) (string, error) {
	if p.TemplateFormat != "f-string" {
		return "", fmt.Errorf("unsupported template format: %s", p.TemplateFormat)
	}

	tmpl, err := template.New("prompt").Parse(p.Template)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, kwargs)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

/*
- Example
	inputVariables := []string{"name", "age"}
	tmpl := "My name is {{.name}} and I am {{.age}} years old."
	pt := prompt.NewPromptTemplate(inputVariables, tmpl, nil)
	formatted, err := pt.Format(map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	})
*/
func NewPromptTemplate(inputVariables []string, tmpl string, partialVariables []string) *PromptTemplate {
	return &PromptTemplate{
		InputVariables:   inputVariables,
		Template:         tmpl,
		TemplateFormat:   "f-string",
		ValidateTemplate: true,
		PartialVariables: partialVariables,
	}
}

func PromptTemplateFromFile(templateFile string, inputVariables []string) (*PromptTemplate, error) {
	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	templateStr := string(templateBytes)
	return NewPromptTemplate(inputVariables, templateStr, nil), nil
}
