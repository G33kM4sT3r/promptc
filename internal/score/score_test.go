package score

import (
	"testing"

	"promptc/internal/prompt"
)

func TestScoreEmpty(t *testing.T) {
	result := Score(prompt.PromptSpec{})
	if result.Total != 0 {
		t.Errorf("empty spec score = %d, want 0", result.Total)
	}
}

func TestScoreFullSpec(t *testing.T) {
	spec := prompt.PromptSpec{
		Role:            "Engineer",
		Objective:       "Explain closures",
		Context:         "For beginners",
		Scope:           []string{"Core concepts"},
		Constraints:     []string{"Simple language"},
		OutputSpec:      []string{"Use headings"},
		QualityCriteria: []string{"Clear"},
	}
	result := Score(spec)
	if result.Total != 100 {
		t.Errorf("full spec score = %d, want 100", result.Total)
	}
}

func TestScorePartial(t *testing.T) {
	spec := prompt.PromptSpec{
		Objective: "Explain closures",
		Scope:     []string{"Core concepts"},
	}
	result := Score(spec)
	if result.Total <= 0 || result.Total >= 100 {
		t.Errorf("partial spec score = %d, want between 1-99", result.Total)
	}
	if result.Breakdown["objective"] != 25 {
		t.Errorf("objective score = %d, want 25", result.Breakdown["objective"])
	}
}
