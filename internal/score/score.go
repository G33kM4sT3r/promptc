package score

import "promptc/internal/prompt"

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
func Score(p prompt.PromptSpec) ScoreResult {
	breakdown := map[string]int{}

	if p.Role != "" {
		breakdown["role"] = 10
	}
	if p.Objective != "" {
		breakdown["objective"] = 25
	}
	if p.Context != "" {
		breakdown["context"] = 15
	}
	if len(p.Scope) > 0 {
		breakdown["scope"] = 15
	}
	if len(p.Constraints) > 0 {
		breakdown["constraints"] = 10
	}
	if len(p.OutputSpec) > 0 {
		breakdown["output"] = 15
	}
	if len(p.QualityCriteria) > 0 {
		breakdown["quality"] = 10
	}

	total := 0
	for _, v := range breakdown {
		total += v
	}

	return ScoreResult{Total: total, Breakdown: breakdown}
}
