package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeDecideRule() rules.Rule {
	return rules.Rule{
		ID: "scope.decide",
		When: func(s slots.Slots) bool {
			return s.Intent == "decide"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.decide_options"),
				t.Get("scope.decide_tradeoffs"),
				t.Get("scope.decide_recommendation"),
			)
		},
	}
}
