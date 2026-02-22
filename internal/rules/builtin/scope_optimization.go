package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeHowToOptimizationRule() rules.Rule {
	return rules.Rule{
		ID: "scope.howto_optimization",
		When: func(s slots.Slots) bool {
			return s.Intent == "howto" && s.Stage == "optimization"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.howto_opt_identify"),
				t.Get("scope.howto_opt_strategies"),
				t.Get("scope.howto_opt_prioritize"),
			)
		},
	}
}
