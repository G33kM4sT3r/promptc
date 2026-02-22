package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeHowToGettingStartedRule() rules.Rule {
	return rules.Rule{
		ID: "scope.howto.getting_started",
		When: func(s slots.Slots) bool {
			return s.Intent == "howto" && s.Stage == "getting-started"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = []string{
				t.Get("scope.howto_gs_goals"),
				t.Get("scope.howto_gs_approach"),
				t.Get("scope.howto_gs_steps"),
			}
		},
	}
}
