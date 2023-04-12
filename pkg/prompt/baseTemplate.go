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

type PromptTemplateBuilder struct {
	InputVariables   []string
	Template         string
	PartialVariables map[string]interface{}
}

type IPromptTemplateBuilder interface {
	AddInputVariables([]string) IPromptTemplateBuilder
	AddTemplate(string) IPromptTemplateBuilder
	AddPartialVariables(map[string]interface{}) IPromptTemplateBuilder
	Build() PromptTemplate
}

func NewPromptTemplateBuilder() *PromptTemplateBuilder {
	return &PromptTemplateBuilder{
		PartialVariables: nil,
		InputVariables:   nil,
	}
}
func (b *PromptTemplateBuilder) AddInputVariables(inputVariables []string) IPromptTemplateBuilder {
	b.InputVariables = inputVariables
	return b
}

func (b *PromptTemplateBuilder) AddTemplate(template string) IPromptTemplateBuilder {
	b.Template = template
	return b
}

func (b *PromptTemplateBuilder) AddPartialVariables(partialVariables map[string]interface{}) IPromptTemplateBuilder {
	b.PartialVariables = partialVariables
	return b
}

func (b *PromptTemplateBuilder) Build() PromptTemplate {
	return PromptTemplate{
		InputVariables:   b.InputVariables,
		Template:         b.Template,
		TemplateFormat:   "f-string",
		ValidateTemplate: true,
		PartialVariables: b.PartialVariables,
	}

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

func PromptTemplateFromFile(templateFile string, inputVariables []string, partialVariables map[string]interface{}) (PromptTemplate, error) {
	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return PromptTemplate{}, err
	}

	templateStr := string(templateBytes)
	return NewPromptTemplateBuilder().AddTemplate(templateStr).AddInputVariables(inputVariables).AddPartialVariables(partialVariables).Build(), nil
}
