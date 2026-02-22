package rules

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/slots"
	"testing"
)

func TestEngineApply(t *testing.T) {
	r1 := Rule{
		ID:   "test.always",
		When: func(_ slots.Slots) bool { return true },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {
			p.Objective = "test"
		},
	}
	r2 := Rule{
		ID:   "test.never",
		When: func(_ slots.Slots) bool { return false },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {
			p.Objective = "should not fire"
		},
	}

	engine := NewEngine([]Rule{r1, r2})
	spec := engine.Apply(slots.Slots{}, nil)

	if spec.Objective != "test" {
		t.Errorf("Objective = %q, want %q", spec.Objective, "test")
	}
}

func TestEngineApplyWithTrace(t *testing.T) {
	r1 := Rule{
		ID:   "test.fires",
		When: func(_ slots.Slots) bool { return true },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {
			p.Objective = "fired"
		},
	}
	r2 := Rule{
		ID:    "test.skipped",
		When:  func(_ slots.Slots) bool { return false },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {},
	}

	engine := NewEngine([]Rule{r1, r2})
	result := engine.ApplyWithTrace(slots.Slots{}, nil)

	if len(result.Applied) != 1 || result.Applied[0] != "test.fires" {
		t.Errorf("Applied = %v, want [test.fires]", result.Applied)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "test.skipped" {
		t.Errorf("Skipped = %v, want [test.skipped]", result.Skipped)
	}
}

func TestEngineOrderDependence(t *testing.T) {
	// Rules are append-only: second rule should append, not override
	r1 := Rule{
		ID:   "scope.first",
		When: func(_ slots.Slots) bool { return true },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {
			p.Scope = append(p.Scope, "first")
		},
	}
	r2 := Rule{
		ID:   "scope.second",
		When: func(_ slots.Slots) bool { return true },
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, _ *i18n.Translator) {
			p.Scope = append(p.Scope, "second")
		},
	}

	engine := NewEngine([]Rule{r1, r2})
	spec := engine.Apply(slots.Slots{}, nil)

	if len(spec.Scope) != 2 || spec.Scope[0] != "first" || spec.Scope[1] != "second" {
		t.Errorf("Scope = %v, want [first, second]", spec.Scope)
	}
}
