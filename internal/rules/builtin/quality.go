package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func QualityBaseRule() rules.Rule {
	return rules.Rule{
		ID: "quality.base",
		When: func(_ slots.Slots) bool {
			return true
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.QualityCriteria = append(p.QualityCriteria,
				t.Get("quality.clear"),
				t.Get("quality.accurate"),
			)
		},
	}
}

func QualityFromIntentRule() rules.Rule {
	return rules.Rule{
		ID: "quality.from_intent",
		When: func(s slots.Slots) bool {
			return s.Intent != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Intent {
			case "generate":
				p.QualityCriteria = append(p.QualityCriteria, t.Get("quality.complete"))
			case "analyze":
				p.QualityCriteria = append(p.QualityCriteria, t.Get("quality.balanced"))
			case "decide":
				p.QualityCriteria = append(p.QualityCriteria, t.Get("quality.fair"))
			}
		},
	}
}
