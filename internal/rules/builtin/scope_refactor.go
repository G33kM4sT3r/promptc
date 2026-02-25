package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeRefactorRule() rules.Rule {
	return rules.Rule{
		ID: "scope.refactor",
		When: func(s slots.Slots) bool {
			return s.Intent == "refactor"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.refactor_problems"),
				t.Get("scope.refactor_approach"),
				t.Get("scope.refactor_before_after"),
				t.Get("scope.refactor_risks"),
			)
		},
	}
}
