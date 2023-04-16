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
	inputVariables   []string
	template         string
	partialVariables map[string]interface{}
}

type promptTemplateBuilder interface {
	AddInputVariables([]string) promptTemplateBuilder
	AddTemplate(string) promptTemplateBuilder
	AddPartialVariables(map[string]interface{}) promptTemplateBuilder
	Build() PromptTemplate
}

func NewPromptTemplateBuilder() *PromptTemplateBuilder {
	return &PromptTemplateBuilder{
		partialVariables: nil,
		inputVariables:   nil,
	}
}
func (b *PromptTemplateBuilder) AddInputVariables(inputVariables []string) promptTemplateBuilder {
	b.inputVariables = inputVariables
	return b
}

func (b *PromptTemplateBuilder) AddTemplate(template string) promptTemplateBuilder {
	b.template = template
	return b
}

func (b *PromptTemplateBuilder) AddPartialVariables(partialVariables map[string]interface{}) promptTemplateBuilder {
	b.partialVariables = partialVariables
	return b
}

func (b *PromptTemplateBuilder) Build() PromptTemplate {
	return PromptTemplate{
		InputVariables:   b.inputVariables,
		Template:         b.template,
		TemplateFormat:   "f-string",
		ValidateTemplate: true,
		PartialVariables: b.partialVariables,
	}

}

func (p *PromptTemplate) Format(args []string, kwargs map[string]interface{}) (string, error) {
	if kwargs == nil {
		kwargs = make(map[string]interface{})
	}
	if p.TemplateFormat != "f-string" {
		return "", fmt.Errorf("unsupported template format: %s", p.TemplateFormat)
	}
	if args != nil {
		for i, v := range p.InputVariables {
			if i < len(args) {
				kwargs[v] = args[i]
				continue
			}
			kwargs[v] = ""
		}
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
