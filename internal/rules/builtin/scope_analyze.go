package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeAnalyzeRule() rules.Rule {
	return rules.Rule{
		ID: "scope.analyze",
		When: func(s slots.Slots) bool {
			return s.Intent == "analyze"
		},
		Apply: func(p *prompt.PromptSpec, _ slots.Slots, t *i18n.Translator) {
			p.Scope = append(p.Scope,
				t.Get("scope.analyze_strengths"),
				t.Get("scope.analyze_evidence"),
				t.Get("scope.analyze_improvements"),
			)
		},
	}
}
