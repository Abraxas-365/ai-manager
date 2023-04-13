package llmmath

import "github.com/Abraxas-365/ai-manager/pkg/prompt"

const Prompt_Template = `Translate a math problem into Python code that can be executed in Python 3 REPL. Use the output of running this code to answer the question.

Question: ${{Question with math problem}}
` + "```python" + `
${{Code that solves the problem and prints the solution}}
` + "```" + `
` + "```output" + `
${{Output of running the code}}
` + "```" + `
Answer: ${{Answer}}

Begin.

Question: What is 37593 * 67?

` + "```python" + `
print(37593 * 67)
` + "```" + `
` + "```output" + `
2518731
` + "```" + `
Answer: 2518731

Question: {{.question}}
`

var Prompt = prompt.NewPromptTemplateBuilder().AddTemplate(Prompt_Template).AddInputVariables([]string{"question"})
