package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeFromExampleObjectRule() rules.Rule {
	return rules.Rule{
		ID: "scope.from_example_object",
		When: func(s slots.Slots) bool {
			for _, e := range s.Entities {
				if e.Role == "example_object" {
					return true
				}
			}
			return false
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			for _, e := range s.Entities {
				if e.Role == "example_object" {
					p.Scope = append(p.Scope, t.Getf("scope.example_reference", e.Text))
				}
			}
		},
	}
}

func ScopeFallbackRule() rules.Rule {
	return rules.Rule{
		ID: "scope.fallback",
		When: func(s slots.Slots) bool {
			return s.Intent != "" && s.Topic != ""
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			if len(p.Scope) == 0 {
				p.Scope = []string{
					t.Get("scope.fallback"),
				}
			}
		},
	}
}
