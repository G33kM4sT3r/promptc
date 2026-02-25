package score

import (
	"strings"

	"promptc/internal/prompt"
)

// ScoreResult holds the total score and per-section breakdown.
type ScoreResult struct {
	Total     int            `json:"total"`
	Breakdown map[string]int `json:"breakdown"`
}

// MaxWeights returns the maximum possible score per section.
func MaxWeights() map[string]int {
	return map[string]int{
		"role":        10,
		"objective":   25,
		"context":     15,
		"scope":       15,
		"constraints": 10,
		"output":      15,
		"quality":     10,
	}
}

// Score evaluates a PromptSpec on completeness, returning 0-100.
// List sections use partial credit based on item count.
func Score(p prompt.PromptSpec) ScoreResult {
	breakdown := map[string]int{}

	if p.Role != "" {
		breakdown["role"] = 10
	}
	if p.Objective != "" {
		breakdown["objective"] = 25
	}

	// Context: 5 per line, capped at 15
	if p.Context != "" {
		lines := strings.Count(p.Context, "\n") + 1
		breakdown["context"] = min(lines*5, 15)
	}

	// Scope: 5 per item, capped at 15
	if len(p.Scope) > 0 {
		breakdown["scope"] = min(len(p.Scope)*5, 15)
	}

	// Constraints: 3 per item, capped at 10
	if len(p.Constraints) > 0 {
		breakdown["constraints"] = min(len(p.Constraints)*3, 10)
	}

	// Output: 5 per item, capped at 15
	if len(p.OutputSpec) > 0 {
		breakdown["output"] = min(len(p.OutputSpec)*5, 15)
	}

	// Quality: 3 per item, capped at 10
	if len(p.QualityCriteria) > 0 {
		breakdown["quality"] = min(len(p.QualityCriteria)*3, 10)
	}

	total := 0
	for _, v := range breakdown {
		total += v
	}

	return ScoreResult{Total: total, Breakdown: breakdown}
}
