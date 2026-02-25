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
		Role:      "Engineer",
		Objective: "Explain closures",
		Context:   "For beginners\nWith examples\nProgressive",
		Scope:     []string{"Core concepts", "Examples", "Edge cases"},
		Constraints: []string{
			"Simple language",
			"No jargon",
			"Define terms",
			"Accessible",
		},
		OutputSpec:      []string{"Use headings", "Summaries", "Examples"},
		QualityCriteria: []string{"Clear", "Accurate", "Accessible", "Thorough"},
	}
	result := Score(spec)
	if result.Total != 100 {
		t.Errorf("rich spec score = %d, want 100", result.Total)
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
}

func TestScorePartialCredit(t *testing.T) {
	tests := []struct {
		name    string
		spec    prompt.PromptSpec
		section string
		wantPts int
	}{
		{"scope_1_item", prompt.PromptSpec{Scope: []string{"one"}}, "scope", 5},
		{"scope_2_items", prompt.PromptSpec{Scope: []string{"one", "two"}}, "scope", 10},
		{"scope_3_items", prompt.PromptSpec{Scope: []string{"one", "two", "three"}}, "scope", 15},
		{"scope_4_items_capped", prompt.PromptSpec{Scope: []string{"a", "b", "c", "d"}}, "scope", 15},
		{"constraints_1", prompt.PromptSpec{Constraints: []string{"one"}}, "constraints", 3},
		{"constraints_2", prompt.PromptSpec{Constraints: []string{"a", "b"}}, "constraints", 6},
		{"constraints_3", prompt.PromptSpec{Constraints: []string{"a", "b", "c"}}, "constraints", 9},
		{"constraints_4_capped", prompt.PromptSpec{Constraints: []string{"a", "b", "c", "d"}}, "constraints", 10},
		{"output_1", prompt.PromptSpec{OutputSpec: []string{"one"}}, "output", 5},
		{"output_3_capped", prompt.PromptSpec{OutputSpec: []string{"a", "b", "c"}}, "output", 15},
		{"quality_1", prompt.PromptSpec{QualityCriteria: []string{"one"}}, "quality", 3},
		{"quality_4_capped", prompt.PromptSpec{QualityCriteria: []string{"a", "b", "c", "d"}}, "quality", 10},
		{"context_1_line", prompt.PromptSpec{Context: "single line"}, "context", 5},
		{"context_2_lines", prompt.PromptSpec{Context: "line1\nline2"}, "context", 10},
		{"context_3_lines_capped", prompt.PromptSpec{Context: "a\nb\nc"}, "context", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Score(tt.spec)
			got := result.Breakdown[tt.section]
			if got != tt.wantPts {
				t.Errorf("%s = %d, want %d", tt.section, got, tt.wantPts)
			}
		})
	}
}

func TestScoreObjectiveAndRole(t *testing.T) {
	spec := prompt.PromptSpec{
		Objective: "Explain closures",
		Role:      "instructor",
	}
	result := Score(spec)
	if result.Breakdown["objective"] != 25 {
		t.Errorf("objective = %d, want 25", result.Breakdown["objective"])
	}
	if result.Breakdown["role"] != 10 {
		t.Errorf("role = %d, want 10", result.Breakdown["role"])
	}
}
