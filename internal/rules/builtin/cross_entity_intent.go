package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func CrossEntityIntentRule() rules.Rule {
	return rules.Rule{
		ID: "cross.entity_intent",
		When: func(s slots.Slots) bool {
			return len(s.Entities) > 0 && s.Intent != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			if s.Intent == "decide" && len(s.Entities) >= 2 {
				p.Scope = append(p.Scope, t.Get("cross.multi_entity_decide_scope"))
			}

			for _, e := range s.Entities {
				if e.Role != "implementation_medium" {
					continue
				}
				switch s.Intent {
				case "generate":
					p.Constraints = append(p.Constraints, t.Getf("cross.medium_generate_constraint", e.Text))
				case "debug":
					p.Scope = append(p.Scope, t.Getf("cross.medium_debug_scope", e.Text))
				case "explain":
					p.Scope = append(p.Scope, t.Getf("cross.medium_explain_scope", e.Text))
				}
				break // only first implementation_medium
			}
		},
	}
}
