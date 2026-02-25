package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func CrossStyleAudienceRule() rules.Rule {
	return rules.Rule{
		ID: "cross.style_audience",
		When: func(s slots.Slots) bool {
			return s.Audience == "beginners" && s.Style == "technical"
		},
		// Must fire after ConstraintsFromStyleRule which adds constraints.technical.
		// Filters by translated string equality — both rules use the same Translator
		// so the strings will always match within a given language.
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			technicalStr := t.Get("constraints.technical")
			replacement := t.Get("cross.beginner_technical_constraint")

			filtered := make([]string, 0, len(p.Constraints))
			for _, c := range p.Constraints {
				if c != technicalStr {
					filtered = append(filtered, c)
				}
			}
			filtered = append(filtered, replacement)
			p.Constraints = filtered
		},
	}
}
