package render

import (
	"testing"

	"promptc/internal/prompt"

	"gopkg.in/yaml.v3"
)

func TestYAMLRenderer(t *testing.T) {
	spec := prompt.PromptSpec{
		Role:      "Software engineer",
		Objective: "Explain closures",
		Scope:     []string{"Core concepts"},
	}

	r := &YAMLRenderer{}
	output := r.Render(spec)

	var parsed prompt.PromptSpec
	if err := yaml.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("invalid YAML output: %v", err)
	}
	if parsed.Objective != "Explain closures" {
		t.Errorf("objective = %q, want %q", parsed.Objective, "Explain closures")
	}
}

func TestYAMLRendererEmpty(t *testing.T) {
	r := &YAMLRenderer{}
	output := r.Render(prompt.PromptSpec{})
	if output != "{}\n" {
		t.Errorf("empty spec YAML = %q, want %q", output, "{}\n")
	}
}
