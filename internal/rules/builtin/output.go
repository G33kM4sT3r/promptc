package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func OutputFromFormatRule() rules.Rule {
	return rules.Rule{
		ID: "output.from_format",
		When: func(s slots.Slots) bool {
			return s.Format != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Format {
			case "bullets":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.bullets"))
			case "steps":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.steps"))
			case "prose":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.prose"))
			case "table":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.table"))
			case "code":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.code"))
			}
		},
	}
}

func OutputFromIntentRule() rules.Rule {
	return rules.Rule{
		ID: "output.from_intent",
		When: func(s slots.Slots) bool {
			return s.Intent != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Intent {
			case "howto":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.howto"))
			case "explain":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.explain"))
			case "generate":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.generate"))
			case "analyze":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.analyze"))
			case "decide":
				p.OutputSpec = append(p.OutputSpec, t.Get("output.decide"))
			}
		},
	}
}

func OutputFromDepthRule() rules.Rule {
	return rules.Rule{
		ID: "output.from_depth",
		When: func(s slots.Slots) bool {
			return s.Depth != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			if s.Depth == "short" {
				p.OutputSpec = append(p.OutputSpec,
					t.Get("output.brief"),
				)
			}
		},
	}
}

func OutputFallbackRule() rules.Rule {
	return rules.Rule{
		ID: "output.fallback",
		When: func(_ slots.Slots) bool {
			return true
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			if len(p.OutputSpec) == 0 {
				p.OutputSpec = append(p.OutputSpec,
					t.Get("output.fallback"),
				)
			}
		},
	}
}
