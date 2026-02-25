package builtin

import (
	"strings"
	"testing"

	"promptc/internal/prompt"
	"promptc/internal/slots"
)

// ---------------------------------------------------------------------------
// CrossAudienceIntentRule
// ---------------------------------------------------------------------------

func TestCrossAudienceIntentRule(t *testing.T) {
	rule := CrossAudienceIntentRule()
	tr := loadTestTranslator(t)

	t.Run("When_audience_and_intent", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "beginners", Intent: "explain"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_no_audience", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_beginner_explain", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "beginners", Intent: "explain"}, tr)
		if len(spec.Scope) != 1 {
			t.Fatalf("len(Scope) = %d, want 1", len(spec.Scope))
		}
	})

	t.Run("Apply_advanced_generate", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "advanced", Intent: "generate"}, tr)
		if len(spec.Constraints) != 1 {
			t.Fatalf("len(Constraints) = %d, want 1", len(spec.Constraints))
		}
	})

	t.Run("Apply_no_match_no_change", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "beginners", Intent: "summarize"}, tr)
		if len(spec.Scope) != 0 || len(spec.Constraints) != 0 {
			t.Error("expected no changes for unmatched combo")
		}
	})
}

// ---------------------------------------------------------------------------
// CrossEntityIntentRule
// ---------------------------------------------------------------------------

func TestCrossEntityIntentRule(t *testing.T) {
	rule := CrossEntityIntentRule()
	tr := loadTestTranslator(t)

	t.Run("When_entity_and_intent", func(t *testing.T) {
		s := slots.Slots{
			Intent:   "generate",
			Entities: []slots.Entity{{Text: "Python", Role: "implementation_medium"}},
		}
		if !rule.When(s) {
			t.Error("expected true")
		}
	})

	t.Run("When_no_entities", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "generate"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_medium_generate", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{
			Intent:   "generate",
			Entities: []slots.Entity{{Text: "Python", Role: "implementation_medium"}},
		}
		rule.Apply(&spec, s, tr)
		if len(spec.Constraints) != 1 {
			t.Fatalf("len(Constraints) = %d, want 1", len(spec.Constraints))
		}
		if !strings.Contains(spec.Constraints[0], "Python") {
			t.Errorf("expected constraint to contain 'Python', got %q", spec.Constraints[0])
		}
	})

	t.Run("Apply_multi_entity_decide", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{
			Intent: "decide",
			Entities: []slots.Entity{
				{Text: "React", Role: "implementation_medium"},
				{Text: "Vue", Role: "implementation_medium"},
			},
		}
		rule.Apply(&spec, s, tr)
		if len(spec.Scope) != 1 {
			t.Fatalf("len(Scope) = %d, want 1", len(spec.Scope))
		}
	})
}

// ---------------------------------------------------------------------------
// CrossStageDepthRule
// ---------------------------------------------------------------------------

func TestCrossStageDepthRule(t *testing.T) {
	rule := CrossStageDepthRule()
	tr := loadTestTranslator(t)

	t.Run("When_stage_and_depth", func(t *testing.T) {
		if !rule.When(slots.Slots{Stage: "getting-started", Depth: "deep"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_no_stage", func(t *testing.T) {
		if rule.When(slots.Slots{Depth: "deep"}) {
			t.Error("expected false")
		}
	})

	t.Run("When_no_depth", func(t *testing.T) {
		if rule.When(slots.Slots{Stage: "getting-started"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_gs_deep", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Stage: "getting-started", Depth: "deep"}, tr)
		if spec.Context == "" {
			t.Error("expected context to be set")
		}
	})

	t.Run("Apply_impl_short", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Stage: "implementation", Depth: "short"}, tr)
		if len(spec.Constraints) != 1 {
			t.Fatalf("len(Constraints) = %d, want 1", len(spec.Constraints))
		}
	})

	t.Run("Apply_opt_deep", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Stage: "optimization", Depth: "deep"}, tr)
		if len(spec.Scope) != 1 {
			t.Fatalf("len(Scope) = %d, want 1", len(spec.Scope))
		}
	})
}

// ---------------------------------------------------------------------------
// CrossAudienceDepthRule
// ---------------------------------------------------------------------------

func TestCrossAudienceDepthRule(t *testing.T) {
	rule := CrossAudienceDepthRule()
	tr := loadTestTranslator(t)

	t.Run("When_beginner_deep", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "beginners", Depth: "deep"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_advanced_short", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "advanced", Depth: "short"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_beginner_standard_false", func(t *testing.T) {
		if rule.When(slots.Slots{Audience: "beginners", Depth: "standard"}) {
			t.Error("expected false for non-conflict combo")
		}
	})

	t.Run("Apply_beginner_deep_adds_context", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "beginners", Depth: "deep"}, tr)
		if spec.Context == "" {
			t.Error("expected context for beginner+deep")
		}
	})

	t.Run("Apply_advanced_short_adds_constraint", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "advanced", Depth: "short"}, tr)
		if len(spec.Constraints) != 1 {
			t.Fatalf("len(Constraints) = %d, want 1", len(spec.Constraints))
		}
	})
}

// ---------------------------------------------------------------------------
// CrossStyleAudienceRule
// ---------------------------------------------------------------------------

func TestCrossStyleAudienceRule(t *testing.T) {
	rule := CrossStyleAudienceRule()
	tr := loadTestTranslator(t)

	t.Run("When_beginner_technical", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "beginners", Style: "technical"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_no_conflict", func(t *testing.T) {
		if rule.When(slots.Slots{Audience: "advanced", Style: "technical"}) {
			t.Error("expected false for non-conflict")
		}
	})

	t.Run("Apply_replaces_technical_constraint", func(t *testing.T) {
		spec := prompt.PromptSpec{
			Constraints: []string{"Use precise technical terminology"},
		}
		rule.Apply(&spec, slots.Slots{Audience: "beginners", Style: "technical"}, tr)
		for _, c := range spec.Constraints {
			if c == "Use precise technical terminology" {
				t.Error("expected technical constraint to be replaced")
			}
		}
		found := false
		for _, c := range spec.Constraints {
			if c == "Define technical terms before using them" {
				found = true
			}
		}
		if !found {
			t.Error("expected replacement constraint")
		}
	})
}
