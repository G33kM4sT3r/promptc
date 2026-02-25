package builtin

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

func CrossStageDepthRule() rules.Rule {
	return rules.Rule{
		ID: "cross.stage_depth",
		When: func(s slots.Slots) bool {
			return s.Stage != "" && s.Depth != ""
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			switch {
			case s.Stage == "getting-started" && s.Depth == "deep":
				appendContext(p, []string{t.Get("cross.gs_deep_context")})
			case s.Stage == "implementation" && s.Depth == "short":
				p.Constraints = append(p.Constraints, t.Get("cross.impl_short_constraint"))
			case s.Stage == "optimization" && s.Depth == "deep":
				p.Scope = append(p.Scope, t.Get("cross.opt_deep_scope"))
			}
		},
	}
}
