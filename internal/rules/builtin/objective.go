package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ObjectiveRule() rules.Rule {
	return rules.Rule{
		ID: "objective.base",
		When: func(s slots.Slots) bool {
			return s.Intent != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			topic := s.Topic
			if topic == "" {
				topic = t.Get("objective.fallback_topic")
			}

			switch s.Intent {
			case "explain":
				p.Objective = t.Getf("objective.explain", topic)
			case "howto":
				p.Objective = t.Getf("objective.howto", topic)
			case "generate":
				p.Objective = t.Getf("objective.generate", topic)
			case "analyze":
				p.Objective = t.Getf("objective.analyze", topic)
			case "decide":
				p.Objective = t.Getf("objective.decide", topic)
			default:
				p.Objective = t.Get("objective.fallback")
			}
		},
	}
}
