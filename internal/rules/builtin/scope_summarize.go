package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeSummarizeRule() rules.Rule {
	return rules.Rule{
		ID: "scope.summarize",
		When: func(s slots.Slots) bool {
			return s.Intent == "summarize"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.summarize_key_points"),
				t.Get("scope.summarize_relationships"),
				t.Get("scope.summarize_takeaways"),
			)
		},
	}
}
