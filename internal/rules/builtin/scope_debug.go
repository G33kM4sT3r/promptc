package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeDebugRule() rules.Rule {
	return rules.Rule{
		ID: "scope.debug",
		When: func(s slots.Slots) bool {
			return s.Intent == "debug"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.debug_symptoms"),
				t.Get("scope.debug_hypotheses"),
				t.Get("scope.debug_investigation"),
				t.Get("scope.debug_resolution"),
			)
		},
	}
}
