package render

import (
	"encoding/json"

	"promptc/internal/prompt"
)

// JSONRenderer renders a PromptSpec as indented JSON.
type JSONRenderer struct{}

// Render marshals the PromptSpec to pretty-printed JSON.
func (r *JSONRenderer) Render(p prompt.PromptSpec) string {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

// RenderScore marshals a ScoreBreakdown to pretty-printed JSON.
func (r *JSONRenderer) RenderScore(sb ScoreBreakdown) string {
	b, err := json.MarshalIndent(sb, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
