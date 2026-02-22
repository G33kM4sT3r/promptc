package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeGenerateRule() rules.Rule {
	return rules.Rule{
		ID: "scope.generate",
		When: func(s slots.Slots) bool {
			return s.Intent == "generate"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.generate_complete"),
				t.Get("scope.generate_conventions"),
				t.Get("scope.generate_structure"),
			)
		},
	}
}
