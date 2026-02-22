package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func ScopeExplainRule() rules.Rule {
	return rules.Rule{
		ID: "scope.explain",
		When: func(s slots.Slots) bool {
			return s.Intent == "explain"
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch s.Depth {
			case "short":
				p.Scope = []string{
					t.Get("scope.explain_short_overview"),
					t.Get("scope.explain_short_key"),
				}
			case "deep":
				p.Scope = []string{
					t.Get("scope.explain_deep_detail"),
					t.Get("scope.explain_deep_examples"),
					t.Get("scope.explain_deep_nuances"),
				}
			default:
				p.Scope = []string{
					t.Get("scope.explain_main"),
					t.Get("scope.explain_aspects"),
				}
			}
		},
	}
}
