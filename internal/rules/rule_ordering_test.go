package rules_test

import (
	"testing"

	"promptc/internal/config"
	"promptc/internal/rules"
	"promptc/internal/rules/builtin"
)

// canonicalRuleOrder returns the rules in the same order as pipeline.newEngine().
// Any change to this list must be reflected in pipeline.go and vice versa.
func canonicalRuleOrder() []rules.Rule {
	return []rules.Rule{
		builtin.ObjectiveRule(),
		builtin.RoleFromEntitiesRule(),

		// Context
		builtin.ContextFromStageRule(),
		builtin.ContextFromAudienceRule(),
		builtin.ContextFromTargetObjectRule(),

		// Scope (specific rules before fallback)
		builtin.ScopeHowToGettingStartedRule(),
		builtin.ScopeHowToImplementationRule(),
		builtin.ScopeHowToOptimizationRule(),
		builtin.ScopeExplainRule(),
		builtin.ScopeGenerateRule(),
		builtin.ScopeAnalyzeRule(),
		builtin.ScopeDecideRule(),
		builtin.ScopeDebugRule(),
		builtin.ScopeRefactorRule(),
		builtin.ScopeSummarizeRule(),
		builtin.ScopeFromExampleObjectRule(),
		builtin.ScopeFallbackRule(),

		// Constraints
		builtin.ConstraintsFromDepthRule(),
		builtin.ConstraintsFromAudienceRule(),
		builtin.ConstraintsFromStyleRule(),
		builtin.ConstraintsFromEntitiesRule(),
		builtin.ConstraintsFromConstraintObjectRule(),

		// Output
		builtin.OutputFromFormatRule(),
		builtin.OutputFromIntentRule(),
		builtin.OutputFromDepthRule(),
		builtin.OutputFallbackRule(),

		// Quality
		builtin.QualityBaseRule(),
		builtin.QualityFromIntentRule(),

		// Tier enrichment
		builtin.EnrichFromTierRule(config.EnrichmentsConfig{}),

		// Cross-field interactions
		builtin.CrossAudienceIntentRule(),
		builtin.CrossEntityIntentRule(),
		builtin.CrossStageDepthRule(),
		builtin.CrossAudienceDepthRule(),
		builtin.CrossStyleAudienceRule(),
	}
}

func TestRuleOrderingConstraints(t *testing.T) {
	ruleList := canonicalRuleOrder()

	// Build position index: rule ID → position
	position := make(map[string]int, len(ruleList))
	for i, r := range ruleList {
		position[r.ID] = i
	}

	// Verify we have exactly 34 rules
	if len(ruleList) != 34 {
		t.Fatalf("expected 34 canonical rules, got %d", len(ruleList))
	}

	// Define ordering constraints: {before, after, reason}
	type constraint struct {
		before string
		after  string
		reason string
	}

	constraints := []constraint{
		// Context ordering: stage sets context before audience and target append
		{"context.from_stage", "context.from_audience", "audience appends to stage-set context"},
		{"context.from_stage", "context.from_target_object", "target appends to stage-set context"},

		// Scope: all specific rules before fallback
		{"scope.howto.getting_started", "scope.fallback", "specific scope before fallback"},
		{"scope.howto.implementation", "scope.fallback", "specific scope before fallback"},
		{"scope.howto_optimization", "scope.fallback", "specific scope before fallback"},
		{"scope.explain", "scope.fallback", "specific scope before fallback"},
		{"scope.generate", "scope.fallback", "specific scope before fallback"},
		{"scope.analyze", "scope.fallback", "specific scope before fallback"},
		{"scope.decide", "scope.fallback", "specific scope before fallback"},
		{"scope.debug", "scope.fallback", "specific scope before fallback"},
		{"scope.refactor", "scope.fallback", "specific scope before fallback"},
		{"scope.summarize", "scope.fallback", "specific scope before fallback"},

		// Scope: intent scope sets before example appends
		{"scope.explain", "scope.from_example_object", "intent scope sets before example appends"},
		{"scope.generate", "scope.from_example_object", "intent scope sets before example appends"},

		// Output: specific rules before fallback
		{"output.from_format", "output.fallback", "format before fallback"},
		{"output.from_intent", "output.fallback", "intent before fallback"},

		// Quality: base before intent-specific
		{"quality.base", "quality.from_intent", "base quality before intent-specific quality"},

		// Tier enrichment: after all base rules, before cross-field
		{"quality.from_intent", "enrich.from_tier", "tier enrichment after all base rules"},

		// Cross-field: after tier enrichment
		{"enrich.from_tier", "cross.audience_intent", "cross-field after tier enrichment"},
		{"enrich.from_tier", "cross.entity_intent", "cross-field after tier enrichment"},
		{"enrich.from_tier", "cross.stage_depth", "cross-field after tier enrichment"},
		{"enrich.from_tier", "cross.audience_depth", "cross-field after tier enrichment"},
		{"enrich.from_tier", "cross.style_audience", "cross-field after tier enrichment"},

		// Style guard must be last cross-field (modifies constraints set by others)
		{"cross.audience_depth", "cross.style_audience", "style guard after depth conflict resolution"},
	}

	for _, c := range constraints {
		posBefore, okBefore := position[c.before]
		posAfter, okAfter := position[c.after]

		if !okBefore {
			t.Errorf("constraint references unknown rule %q", c.before)
			continue
		}
		if !okAfter {
			t.Errorf("constraint references unknown rule %q", c.after)
			continue
		}

		if posBefore >= posAfter {
			t.Errorf("ordering violation: %q (pos %d) must come before %q (pos %d): %s",
				c.before, posBefore, c.after, posAfter, c.reason)
		}
	}
}
