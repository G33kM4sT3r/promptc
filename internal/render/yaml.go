package render

import (
	"promptc/internal/prompt"

	"gopkg.in/yaml.v3"
)

// YAMLRenderer renders a PromptSpec as YAML.
type YAMLRenderer struct{}

// Render marshals the PromptSpec to YAML.
func (r *YAMLRenderer) Render(p prompt.PromptSpec) string {
	b, err := yaml.Marshal(p)
	if err != nil {
		return "{}\n"
	}
	return string(b)
}

// RenderScore marshals a ScoreBreakdown to YAML.
func (r *YAMLRenderer) RenderScore(sb ScoreBreakdown) string {
	b, err := yaml.Marshal(sb)
	if err != nil {
		return "{}\n"
	}
	return string(b)
}
