package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func RoleFromEntitiesRule() rules.Rule {
	return rules.Rule{
		ID: "role.from_entity",
		When: func(s slots.Slots) bool {
			return len(s.Entities) > 0
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			for _, e := range s.Entities {
				if e.Role == "implementation_medium" {
					p.Role = t.Getf("role.from_entity", e.Text)
					return
				}
			}
		},
	}
}
