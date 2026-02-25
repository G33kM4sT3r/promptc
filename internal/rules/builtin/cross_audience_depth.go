package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func CrossAudienceDepthRule() rules.Rule {
	return rules.Rule{
		ID: "cross.audience_depth",
		When: func(s slots.Slots) bool {
			return (s.Audience == "beginners" && s.Depth == "deep") ||
				(s.Audience == "advanced" && s.Depth == "short")
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch {
			case s.Audience == "beginners" && s.Depth == "deep":
				appendContext(p, []string{t.Get("cross.beginner_deep_context")})
			case s.Audience == "advanced" && s.Depth == "short":
				p.Constraints = append(p.Constraints, t.Get("cross.advanced_short_constraint"))
			}
		},
	}
}
