package builtin

import (
	"promptc/internal/config"
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/rules"
	"promptc/internal/slots"
)

// EnrichFromTierRule returns a rule that enriches the PromptSpec based on
// the calculated tier and intent, using content from the enrichments config.
func EnrichFromTierRule(e config.EnrichmentsConfig) rules.Rule {
	return rules.Rule{
		ID: "enrich.from_tier",
		When: func(s slots.Slots) bool {
			return s.Intent != "" && s.Tier != "" && s.Tier != "minimal"
		},
		Apply: func(p *prompt.PromptSpec, s slots.Slots, t *i18n.Translator) {
			intent, tier := s.Intent, s.Tier

			if p.Role == "" {
				if key := lookupRole(e.Roles, intent, tier); key != "" {
					p.Role = t.Get(key)
				}
			}

			appendContext(p, translateKeys(t, lookupKeys(e.Context, intent, tier)))
			p.Scope = append(p.Scope, translateKeys(t, lookupKeys(e.Scope, intent, tier))...)
			p.Constraints = append(p.Constraints, translateKeys(t, lookupKeys(e.Constraints, intent, tier))...)
			p.OutputSpec = append(p.OutputSpec, translateKeys(t, lookupKeys(e.Output, intent, tier))...)
			p.QualityCriteria = append(p.QualityCriteria, translateKeys(t, lookupKeys(e.Quality, intent, tier))...)
		},
	}
}

// lookupKeys returns translation keys for a given intent+tier from a nested map.
func lookupKeys(m map[string]map[string][]string, intent, tier string) []string {
	if m == nil {
		return nil
	}
	if byTier, ok := m[intent]; ok {
		return byTier[tier]
	}
	return nil
}

// lookupRole returns a single translation key for a given intent+tier.
func lookupRole(m map[string]map[string]string, intent, tier string) string {
	if m == nil {
		return ""
	}
	if byTier, ok := m[intent]; ok {
		return byTier[tier]
	}
	return ""
}

// translateKeys resolves a slice of translation keys into translated strings.
func translateKeys(t *i18n.Translator, keys []string) []string {
	if len(keys) == 0 {
		return nil
	}
	out := make([]string, len(keys))
	for i, key := range keys {
		out[i] = t.Get(key)
	}
	return out
}

// appendContext appends lines to the Context field with newline separators.
func appendContext(p *prompt.PromptSpec, lines []string) {
	for _, line := range lines {
		if p.Context == "" {
			p.Context = line
		} else {
			p.Context += "\n" + line
		}
	}
}
