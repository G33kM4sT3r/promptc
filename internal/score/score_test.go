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

func TestMaxWeights(t *testing.T) {
	weights := MaxWeights()
	total := 0
	for _, v := range weights {
		total += v
	}
	if total != 100 {
		t.Errorf("MaxWeights sum = %d, want 100", total)
	}
	if weights["objective"] != 25 {
		t.Errorf("objective weight = %d, want 25", weights["objective"])
	}
	if weights["role"] != 10 {
		t.Errorf("role weight = %d, want 10", weights["role"])
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
