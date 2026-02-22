package rules

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/slots"
)

// ApplyResult holds the result of rule application with tracing info.
type ApplyResult struct {
	Spec    prompt.PromptSpec
	Applied []string // IDs of rules that fired
	Skipped []string // IDs of rules whose guard returned false
}

type Engine struct {
	rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) Apply(s slots.Slots, t *i18n.Translator) prompt.PromptSpec {
	result := e.ApplyWithTrace(s, t)
	return result.Spec
}

// ApplyWithTrace applies rules and returns tracing information.
func (e *Engine) ApplyWithTrace(s slots.Slots, t *i18n.Translator) ApplyResult {
	var result ApplyResult

	for _, r := range e.rules {
		if r.When(s) {
			r.Apply(&result.Spec, s, t)
			result.Applied = append(result.Applied, r.ID)
		} else {
			result.Skipped = append(result.Skipped, r.ID)
		}
	}

	// Deduplicate slice fields
	result.Spec.Scope = dedup(result.Spec.Scope)
	result.Spec.Constraints = dedup(result.Spec.Constraints)
	result.Spec.OutputSpec = dedup(result.Spec.OutputSpec)
	result.Spec.QualityCriteria = dedup(result.Spec.QualityCriteria)

	return result
}

func dedup(items []string) []string {
	if len(items) == 0 {
		return items
	}
	seen := make(map[string]bool)
	result := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
