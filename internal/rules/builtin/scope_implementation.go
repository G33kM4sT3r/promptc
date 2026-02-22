package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeHowToImplementationRule() rules.Rule {
	return rules.Rule{
		ID: "scope.howto.implementation",
		When: func(s slots.Slots) bool {
			return s.Intent == "howto" && s.Stage == "implementation"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = []string{
				t.Get("scope.howto_impl_components"),
				t.Get("scope.howto_impl_approach"),
				t.Get("scope.howto_impl_pitfalls"),
			}
		},
	}
}
