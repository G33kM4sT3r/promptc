package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func CrossAudienceIntentRule() rules.Rule {
	return rules.Rule{
		ID: "cross.audience_intent",
		When: func(s slots.Slots) bool {
			return s.Audience != "" && s.Intent != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch {
			case s.Audience == "beginners" && s.Intent == "explain":
				p.Scope = append(p.Scope, t.Get("cross.beginner_explain_scope"))
			case s.Audience == "beginners" && s.Intent == "generate":
				p.Constraints = append(p.Constraints, t.Get("cross.beginner_generate_constraint"))
			case s.Audience == "beginners" && s.Intent == "debug":
				p.Scope = append(p.Scope, t.Get("cross.beginner_debug_scope"))
			case s.Audience == "advanced" && s.Intent == "analyze":
				p.Scope = append(p.Scope, t.Get("cross.advanced_analyze_scope"))
			case s.Audience == "advanced" && s.Intent == "generate":
				p.Constraints = append(p.Constraints, t.Get("cross.advanced_generate_constraint"))
			}
		},
	}
}
