package builtin

import (
	"testing"

	"promptc/internal/config"
	"promptc/internal/prompt"
	"promptc/internal/slots"
)

func loadTestEnrichments(t *testing.T) config.EnrichmentsConfig {
	t.Helper()
	dir := findProjectRoot(t)
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg.Enrichments
}

func TestEnrichFromTierRule_When(t *testing.T) {
	rule := EnrichFromTierRule(config.EnrichmentsConfig{})

	t.Run("fires_when_intent_and_tier", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "explain", Tier: "standard"}) {
			t.Error("expected true")
		}
	})

	t.Run("skips_when_no_intent", func(t *testing.T) {
		if rule.When(slots.Slots{Tier: "standard"}) {
			t.Error("expected false when intent empty")
		}
	})

	t.Run("skips_when_no_tier", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false when tier empty")
		}
	})

	t.Run("skips_minimal_tier", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain", Tier: "minimal"}) {
			t.Error("expected false for minimal tier")
		}
	})
}

func TestEnrichFromTierRule_Apply_Role(t *testing.T) {
	e := loadTestEnrichments(t)
	tr := loadTestTranslator(t)
	rule := EnrichFromTierRule(e)

	t.Run("sets_role_when_empty", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain", Tier: "standard"}, tr)
		if spec.Role == "" {
			t.Error("Role should be set by enrichment")
		}
		if spec.Role != "knowledgeable instructor" {
			t.Errorf("Role = %q, want %q", spec.Role, "knowledgeable instructor")
		}
	})

	t.Run("does_not_overwrite_existing_role", func(t *testing.T) {
		spec := prompt.PromptSpec{Role: "You are experienced in working with Python."}
		rule.Apply(&spec, slots.Slots{Intent: "explain", Tier: "standard"}, tr)
		if spec.Role != "You are experienced in working with Python." {
			t.Errorf("Role = %q, should not be overwritten", spec.Role)
		}
	})
}

func TestEnrichFromTierRule_Apply_Scope(t *testing.T) {
	e := loadTestEnrichments(t)
	tr := loadTestTranslator(t)
	rule := EnrichFromTierRule(e)

	t.Run("appends_scope_items", func(t *testing.T) {
		spec := prompt.PromptSpec{Scope: []string{"existing"}}
		rule.Apply(&spec, slots.Slots{Intent: "explain", Tier: "standard"}, tr)
		if len(spec.Scope) <= 1 {
			t.Error("expected enrichment to append scope items")
		}
	})

	t.Run("rich_adds_more_than_standard", func(t *testing.T) {
		var specStd prompt.PromptSpec
		rule.Apply(&specStd, slots.Slots{Intent: "explain", Tier: "standard"}, tr)

		var specRich prompt.PromptSpec
		rule.Apply(&specRich, slots.Slots{Intent: "explain", Tier: "rich"}, tr)

		if len(specRich.Scope) <= len(specStd.Scope) {
			t.Errorf("rich scope (%d items) should be larger than standard (%d items)",
				len(specRich.Scope), len(specStd.Scope))
		}
	})
}

func TestEnrichFromTierRule_Apply_Context(t *testing.T) {
	e := loadTestEnrichments(t)
	tr := loadTestTranslator(t)
	rule := EnrichFromTierRule(e)

	t.Run("appends_context", func(t *testing.T) {
		spec := prompt.PromptSpec{Context: "Existing context."}
		rule.Apply(&spec, slots.Slots{Intent: "explain", Tier: "standard"}, tr)
		if spec.Context == "Existing context." {
			t.Error("expected enrichment to append to context")
		}
	})

	t.Run("sets_context_when_empty", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain", Tier: "standard"}, tr)
		if spec.Context == "" {
			t.Error("expected enrichment to set context")
		}
	})
}

func TestEnrichFromTierRule_Apply_MissingIntent(t *testing.T) {
	e := loadTestEnrichments(t)
	tr := loadTestTranslator(t)
	rule := EnrichFromTierRule(e)

	t.Run("unknown_intent_no_panic", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "unknown_intent", Tier: "standard"}, tr)
		// Should not panic and should not modify spec
	})
}
