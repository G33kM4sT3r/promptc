package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ContextFromStageRule() rules.Rule {
	return rules.Rule{
		ID: "context.from_stage",
		When: func(s slots.Slots) bool {
			return s.Stage != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Stage {
			case "getting-started":
				p.Context = t.Get("context.getting_started")
			case "implementation":
				p.Context = t.Get("context.implementation")
			case "optimization":
				p.Context = t.Get("context.optimization")
			}
		},
	}
}

func ContextFromAudienceRule() rules.Rule {
	return rules.Rule{
		ID: "context.from_audience",
		When: func(s slots.Slots) bool {
			return s.Audience != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			ctx := ""

			switch s.Audience {
			case "beginners":
				ctx = t.Get("context.audience_beginners")
			case "advanced":
				ctx = t.Get("context.audience_advanced")
			}

			if ctx == "" {
				return
			}

			if p.Context == "" {
				p.Context = ctx
			} else {
				p.Context = p.Context + " " + ctx
			}
		},
	}
}

func ContextFromTargetObjectRule() rules.Rule {
	return rules.Rule{
		ID: "context.from_target_object",
		When: func(s slots.Slots) bool {
			for _, e := range s.Entities {
				if e.Role == "target_object" {
					return true
				}
			}
			return false
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			for _, e := range s.Entities {
				if e.Role == "target_object" {
					if p.Context != "" {
						p.Context += " "
					}
					p.Context += t.Getf("context.target_object", e.Text)
				}
			}
		},
	}
}
