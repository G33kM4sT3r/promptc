package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ConstraintsFromDepthRule() rules.Rule {
	return rules.Rule{
		ID: "constraints.from_depth",
		When: func(s slots.Slots) bool {
			return s.Depth != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Depth {
			case "short":
				p.Constraints = append(p.Constraints,
					t.Get("constraints.concise"),
				)
			case "deep":
				p.Constraints = append(p.Constraints,
					t.Get("constraints.detailed"),
				)
			}
		},
	}
}

func ConstraintsFromAudienceRule() rules.Rule {
	return rules.Rule{
		ID: "constraints.from_audience",
		When: func(s slots.Slots) bool {
			return s.Audience != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Audience {
			case "beginners":
				p.Constraints = append(p.Constraints,
					t.Get("constraints.simple_language"),
					t.Get("constraints.no_jargon"),
				)
			case "advanced":
				p.Constraints = append(p.Constraints,
					t.Get("constraints.technical_accuracy"),
				)
			}
		},
	}
}

func ConstraintsFromStyleRule() rules.Rule {
	return rules.Rule{
		ID: "constraints.from_style",
		When: func(s slots.Slots) bool {
			return s.Style != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Style {
			case "formal":
				p.Constraints = append(p.Constraints, t.Get("constraints.formal"))
			case "casual":
				p.Constraints = append(p.Constraints, t.Get("constraints.casual"))
			case "technical":
				p.Constraints = append(p.Constraints, t.Get("constraints.technical"))
			}
		},
	}
}

func ConstraintsFromEntitiesRule() rules.Rule {
	return rules.Rule{
		ID: "constraints.from_entities",
		When: func(s slots.Slots) bool {
			return len(s.Entities) > 0
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			for _, e := range s.Entities {
				if e.Role == "implementation_medium" {
					p.Constraints = append(p.Constraints,
						t.Getf("constraints.practical_usage", e.Text),
					)
				}
			}
		},
	}
}

func ConstraintsFromConstraintObjectRule() rules.Rule {
	return rules.Rule{
		ID: "constraints.from_constraint_object",
		When: func(s slots.Slots) bool {
			for _, e := range s.Entities {
				if e.Role == "constraint_object" {
					return true
				}
			}
			return false
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			for _, e := range s.Entities {
				if e.Role == "constraint_object" {
					p.Constraints = append(p.Constraints, t.Getf("constraints.no_use", e.Text))
				}
			}
		},
	}
}
