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
	PartialVariables map[string]interface{}
}

func (p *PromptTemplate) Format(kwargs map[string]interface{}) (string, error) {
	if p.TemplateFormat != "f-string" {
		return "", fmt.Errorf("unsupported template format: %s", p.TemplateFormat)
	}

	for k, v := range p.PartialVariables {
		kwargs[k] = v
	}

	tmpl, err := template.New("prompt").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(p.Template)
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

func NewPromptTemplate(inputVariables []string, tmpl string, partialVariables map[string]interface{}) *PromptTemplate {
	return &PromptTemplate{
		InputVariables:   inputVariables,
		Template:         tmpl,
		TemplateFormat:   "f-string",
		ValidateTemplate: true,
		PartialVariables: partialVariables,
	}
}

func PromptTemplateFromFile(templateFile string, inputVariables []string, partialVariables map[string]interface{}) (*PromptTemplate, error) {
	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	templateStr := string(templateBytes)
	return NewPromptTemplate(inputVariables, templateStr, partialVariables), nil
}
