package render

import (
	"encoding/json"
	"testing"

	"promptc/internal/prompt"
)

func TestJSONRenderer(t *testing.T) {
	spec := prompt.PromptSpec{
		Role:            "Software engineer",
		Objective:       "Explain closures",
		Context:         "For beginners",
		Scope:           []string{"Core concepts"},
		Constraints:     []string{"Simple language"},
		OutputSpec:      []string{"Use headings"},
		QualityCriteria: []string{"Clear"},
	}

	r := &JSONRenderer{}
	output := r.Render(spec)

	var parsed prompt.PromptSpec
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if parsed.Objective != "Explain closures" {
		t.Errorf("objective = %q, want %q", parsed.Objective, "Explain closures")
	}
	if parsed.Role != "Software engineer" {
		t.Errorf("role = %q, want %q", parsed.Role, "Software engineer")
	}
}

func TestJSONRendererEmpty(t *testing.T) {
	r := &JSONRenderer{}
	output := r.Render(prompt.PromptSpec{})

	var parsed prompt.PromptSpec
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("invalid JSON for empty spec: %v", err)
	}
}

func TestJSONRendererRenderScore(t *testing.T) {
	r := &JSONRenderer{}
	sb := ScoreBreakdown{
		Total:      65,
		Breakdown:  map[string]int{"objective": 25, "scope": 15},
		MaxWeights: map[string]int{"objective": 25, "scope": 15},
	}
	output := r.RenderScore(sb)

	var parsed ScoreBreakdown
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed.Total != 65 {
		t.Errorf("total = %d, want 65", parsed.Total)
	}
}
