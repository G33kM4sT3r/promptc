package builtin

import (
	"os"
	"path/filepath"
	"testing"

	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/slots"
)

func loadTestTranslator(t *testing.T) *i18n.Translator {
	t.Helper()
	dir := findProjectRoot(t)
	tr, err := i18n.Load(filepath.Join(dir, "languages"), "en", "en")
	if err != nil {
		t.Fatalf("failed to load translator: %v", err)
	}
	return tr
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

// ---------------------------------------------------------------------------
// ObjectiveRule
// ---------------------------------------------------------------------------

func TestObjectiveRule(t *testing.T) {
	rule := ObjectiveRule()
	tr := loadTestTranslator(t)

	t.Run("When_intent_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected true for non-empty intent")
		}
	})

	t.Run("When_intent_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false for empty intent")
		}
	})

	applyTests := []struct {
		name   string
		intent string
		topic  string
		want   string
	}{
		{"explain_with_topic", "explain", "closures", "Explain closures."},
		{"howto_with_topic", "howto", "deploy containers", "Provide guidance on how to deploy containers."},
		{"generate_with_topic", "generate", "a REST API", "Generate a REST API."},
		{"analyze_with_topic", "analyze", "performance metrics", "Analyze performance metrics."},
		{"decide_with_topic", "decide", "a database", "Help decide on a database."},
		{"explain_empty_topic", "explain", "", "Explain the given topic."},
		{"unknown_intent", "foobar", "something", "Provide relevant information."},
	}

	for _, tc := range applyTests {
		t.Run("Apply_"+tc.name, func(t *testing.T) {
			var spec prompt.PromptSpec
			rule.Apply(&spec, slots.Slots{Intent: tc.intent, Topic: tc.topic}, tr)
			if spec.Objective != tc.want {
				t.Errorf("Objective = %q, want %q", spec.Objective, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// RoleFromEntitiesRule
// ---------------------------------------------------------------------------

func TestRoleFromEntitiesRule(t *testing.T) {
	rule := RoleFromEntitiesRule()
	tr := loadTestTranslator(t)

	t.Run("When_entities_present", func(t *testing.T) {
		s := slots.Slots{Entities: []slots.Entity{{Text: "Go", Role: "implementation_medium"}}}
		if !rule.When(s) {
			t.Error("expected true when entities present")
		}
	})

	t.Run("When_entities_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false when entities empty")
		}
	})

	t.Run("Apply_implementation_medium", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{Entities: []slots.Entity{{Text: "PHP", Role: "implementation_medium"}}}
		rule.Apply(&spec, s, tr)
		want := "You are experienced in working with PHP."
		if spec.Role != want {
			t.Errorf("Role = %q, want %q", spec.Role, want)
		}
	})

	t.Run("Apply_non_implementation_medium", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{Entities: []slots.Entity{{Text: "users", Role: "target_object"}}}
		rule.Apply(&spec, s, tr)
		if spec.Role != "" {
			t.Errorf("Role = %q, want empty", spec.Role)
		}
	})

	t.Run("Apply_first_implementation_medium_wins", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{Entities: []slots.Entity{
			{Text: "users", Role: "target_object"},
			{Text: "Python", Role: "implementation_medium"},
			{Text: "Go", Role: "implementation_medium"},
		}}
		rule.Apply(&spec, s, tr)
		want := "You are experienced in working with Python."
		if spec.Role != want {
			t.Errorf("Role = %q, want %q", spec.Role, want)
		}
	})
}

// ---------------------------------------------------------------------------
// ContextFromStageRule
// ---------------------------------------------------------------------------

func TestContextFromStageRule(t *testing.T) {
	rule := ContextFromStageRule()
	tr := loadTestTranslator(t)

	t.Run("When_stage_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Stage: "implementation"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_stage_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	stageTests := []struct {
		stage string
		want  string
	}{
		{"getting-started", "The user is at an early stage and is looking for initial guidance."},
		{"implementation", "The user is actively implementing a solution."},
		{"optimization", "The user is looking to improve an existing solution."},
	}

	for _, tc := range stageTests {
		t.Run("Apply_"+tc.stage, func(t *testing.T) {
			var spec prompt.PromptSpec
			rule.Apply(&spec, slots.Slots{Stage: tc.stage}, tr)
			if spec.Context != tc.want {
				t.Errorf("Context = %q, want %q", spec.Context, tc.want)
			}
		})
	}

	t.Run("Apply_unknown_stage_no_context", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Stage: "unknown"}, tr)
		if spec.Context != "" {
			t.Errorf("Context = %q, want empty", spec.Context)
		}
	})
}

// ---------------------------------------------------------------------------
// ContextFromAudienceRule
// ---------------------------------------------------------------------------

func TestContextFromAudienceRule(t *testing.T) {
	rule := ContextFromAudienceRule()
	tr := loadTestTranslator(t)

	t.Run("When_audience_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "beginners"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_audience_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_beginners_empty_context", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "beginners"}, tr)
		want := "The explanation should be beginner-friendly."
		if spec.Context != want {
			t.Errorf("Context = %q, want %q", spec.Context, want)
		}
	})

	t.Run("Apply_advanced_empty_context", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "advanced"}, tr)
		want := "The explanation may assume advanced prior knowledge."
		if spec.Context != want {
			t.Errorf("Context = %q, want %q", spec.Context, want)
		}
	})

	t.Run("Apply_beginners_appends_to_existing_context", func(t *testing.T) {
		spec := prompt.PromptSpec{Context: "Existing context."}
		rule.Apply(&spec, slots.Slots{Audience: "beginners"}, tr)
		want := "Existing context. The explanation should be beginner-friendly."
		if spec.Context != want {
			t.Errorf("Context = %q, want %q", spec.Context, want)
		}
	})

	t.Run("Apply_unknown_audience_no_change", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "intermediate"}, tr)
		if spec.Context != "" {
			t.Errorf("Context = %q, want empty", spec.Context)
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeExplainRule
// ---------------------------------------------------------------------------

func TestScopeExplainRule(t *testing.T) {
	rule := ScopeExplainRule()
	tr := loadTestTranslator(t)

	t.Run("When_explain", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_explain", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "howto"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_short", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain", Depth: "short"}, tr)
		if len(spec.Scope) != 2 {
			t.Fatalf("len(Scope) = %d, want 2", len(spec.Scope))
		}
		if spec.Scope[0] != "Provide a high-level overview" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})

	t.Run("Apply_deep", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain", Depth: "deep"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Explain the core concepts in detail" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})

	t.Run("Apply_default_no_depth", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain"}, tr)
		if len(spec.Scope) != 2 {
			t.Fatalf("len(Scope) = %d, want 2", len(spec.Scope))
		}
		if spec.Scope[0] != "Explain the main concept" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeHowToGettingStartedRule
// ---------------------------------------------------------------------------

func TestScopeHowToGettingStartedRule(t *testing.T) {
	rule := ScopeHowToGettingStartedRule()
	tr := loadTestTranslator(t)

	t.Run("When_howto_getting_started", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "howto", Stage: "getting-started"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_howto_different_stage", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "howto", Stage: "implementation"}) {
			t.Error("expected false")
		}
	})

	t.Run("When_not_howto", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain", Stage: "getting-started"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "howto", Stage: "getting-started"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Clarify the project goals and requirements" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeHowToImplementationRule
// ---------------------------------------------------------------------------

func TestScopeHowToImplementationRule(t *testing.T) {
	rule := ScopeHowToImplementationRule()
	tr := loadTestTranslator(t)

	t.Run("When_howto_implementation", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "howto", Stage: "implementation"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_howto_wrong_stage", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "howto", Stage: "optimization"}) {
			t.Error("expected false")
		}
	})

	t.Run("When_wrong_intent", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain", Stage: "implementation"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "howto", Stage: "implementation"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Define the core components of the solution" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeHowToOptimizationRule
// ---------------------------------------------------------------------------

func TestScopeHowToOptimizationRule(t *testing.T) {
	rule := ScopeHowToOptimizationRule()
	tr := loadTestTranslator(t)

	t.Run("When_howto_optimization", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "howto", Stage: "optimization"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_howto_wrong_stage", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "howto", Stage: "getting-started"}) {
			t.Error("expected false")
		}
	})

	t.Run("When_wrong_intent", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain", Stage: "optimization"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_appends", func(t *testing.T) {
		spec := prompt.PromptSpec{Scope: []string{"existing item"}}
		rule.Apply(&spec, slots.Slots{Intent: "howto", Stage: "optimization"}, tr)
		if len(spec.Scope) != 4 {
			t.Fatalf("len(Scope) = %d, want 4", len(spec.Scope))
		}
		if spec.Scope[0] != "existing item" {
			t.Errorf("Scope[0] = %q, want %q", spec.Scope[0], "existing item")
		}
		if spec.Scope[1] != "Identify areas for improvement" {
			t.Errorf("Scope[1] = %q", spec.Scope[1])
		}
	})

	t.Run("Apply_empty_scope", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "howto", Stage: "optimization"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeGenerateRule
// ---------------------------------------------------------------------------

func TestScopeGenerateRule(t *testing.T) {
	rule := ScopeGenerateRule()
	tr := loadTestTranslator(t)

	t.Run("When_generate", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "generate"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_generate", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "generate"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Produce complete, ready-to-use output" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeAnalyzeRule
// ---------------------------------------------------------------------------

func TestScopeAnalyzeRule(t *testing.T) {
	rule := ScopeAnalyzeRule()
	tr := loadTestTranslator(t)

	t.Run("When_analyze", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "analyze"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_analyze", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "analyze"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Identify strengths and weaknesses" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeDecideRule
// ---------------------------------------------------------------------------

func TestScopeDecideRule(t *testing.T) {
	rule := ScopeDecideRule()
	tr := loadTestTranslator(t)

	t.Run("When_decide", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "decide"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_decide", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "howto"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "decide"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Present available options clearly" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeDebugRule
// ---------------------------------------------------------------------------

func TestScopeDebugRule(t *testing.T) {
	rule := ScopeDebugRule()
	tr := loadTestTranslator(t)

	t.Run("When_debug", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "debug"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_debug", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "debug"}, tr)
		if len(spec.Scope) != 4 {
			t.Fatalf("len(Scope) = %d, want 4", len(spec.Scope))
		}
		if spec.Scope[0] != "Describe the observed symptoms and expected behavior" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeRefactorRule
// ---------------------------------------------------------------------------

func TestScopeRefactorRule(t *testing.T) {
	rule := ScopeRefactorRule()
	tr := loadTestTranslator(t)

	t.Run("When_refactor", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "refactor"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_refactor", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "refactor"}, tr)
		if len(spec.Scope) != 4 {
			t.Fatalf("len(Scope) = %d, want 4", len(spec.Scope))
		}
		if spec.Scope[0] != "Identify current structural problems" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeSummarizeRule
// ---------------------------------------------------------------------------

func TestScopeSummarizeRule(t *testing.T) {
	rule := ScopeSummarizeRule()
	tr := loadTestTranslator(t)

	t.Run("When_summarize", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "summarize"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_not_summarize", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "summarize"}, tr)
		if len(spec.Scope) != 3 {
			t.Fatalf("len(Scope) = %d, want 3", len(spec.Scope))
		}
		if spec.Scope[0] != "Extract the key points" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})
}

// ---------------------------------------------------------------------------
// ScopeFallbackRule
// ---------------------------------------------------------------------------

func TestScopeFallbackRule(t *testing.T) {
	rule := ScopeFallbackRule()
	tr := loadTestTranslator(t)

	t.Run("When_intent_and_topic", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "explain", Topic: "closures"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_no_intent", func(t *testing.T) {
		if rule.When(slots.Slots{Topic: "closures"}) {
			t.Error("expected false when intent empty")
		}
	})

	t.Run("When_no_topic", func(t *testing.T) {
		if rule.When(slots.Slots{Intent: "explain"}) {
			t.Error("expected false when topic empty")
		}
	})

	t.Run("Apply_scope_empty_sets_fallback", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain", Topic: "closures"}, tr)
		if len(spec.Scope) != 1 {
			t.Fatalf("len(Scope) = %d, want 1", len(spec.Scope))
		}
		if spec.Scope[0] != "Address the main topic clearly" {
			t.Errorf("Scope[0] = %q", spec.Scope[0])
		}
	})

	t.Run("Apply_scope_not_empty_no_change", func(t *testing.T) {
		spec := prompt.PromptSpec{Scope: []string{"already set"}}
		rule.Apply(&spec, slots.Slots{Intent: "explain", Topic: "closures"}, tr)
		if len(spec.Scope) != 1 {
			t.Fatalf("len(Scope) = %d, want 1", len(spec.Scope))
		}
		if spec.Scope[0] != "already set" {
			t.Errorf("Scope[0] = %q, want %q", spec.Scope[0], "already set")
		}
	})
}

// ---------------------------------------------------------------------------
// ConstraintsFromDepthRule
// ---------------------------------------------------------------------------

func TestConstraintsFromDepthRule(t *testing.T) {
	rule := ConstraintsFromDepthRule()
	tr := loadTestTranslator(t)

	t.Run("When_depth_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Depth: "short"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_depth_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_short", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "short"}, tr)
		if len(spec.Constraints) != 1 || spec.Constraints[0] != "Keep the response concise" {
			t.Errorf("Constraints = %v", spec.Constraints)
		}
	})

	t.Run("Apply_deep", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "deep"}, tr)
		if len(spec.Constraints) != 1 || spec.Constraints[0] != "Provide detailed explanations" {
			t.Errorf("Constraints = %v", spec.Constraints)
		}
	})

	t.Run("Apply_standard_no_constraint", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "standard"}, tr)
		if len(spec.Constraints) != 0 {
			t.Errorf("Constraints = %v, want empty", spec.Constraints)
		}
	})
}

// ---------------------------------------------------------------------------
// ConstraintsFromAudienceRule
// ---------------------------------------------------------------------------

func TestConstraintsFromAudienceRule(t *testing.T) {
	rule := ConstraintsFromAudienceRule()
	tr := loadTestTranslator(t)

	t.Run("When_audience_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Audience: "beginners"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_audience_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_beginners", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "beginners"}, tr)
		if len(spec.Constraints) != 2 {
			t.Fatalf("len(Constraints) = %d, want 2", len(spec.Constraints))
		}
		if spec.Constraints[0] != "Use simple and clear language" {
			t.Errorf("Constraints[0] = %q", spec.Constraints[0])
		}
		if spec.Constraints[1] != "Avoid unnecessary jargon" {
			t.Errorf("Constraints[1] = %q", spec.Constraints[1])
		}
	})

	t.Run("Apply_advanced", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Audience: "advanced"}, tr)
		if len(spec.Constraints) != 1 || spec.Constraints[0] != "Focus on technical accuracy" {
			t.Errorf("Constraints = %v", spec.Constraints)
		}
	})
}

// ---------------------------------------------------------------------------
// ConstraintsFromStyleRule
// ---------------------------------------------------------------------------

func TestConstraintsFromStyleRule(t *testing.T) {
	rule := ConstraintsFromStyleRule()
	tr := loadTestTranslator(t)

	t.Run("When_style_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Style: "formal"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_style_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	styleTests := []struct {
		style string
		want  string
	}{
		{"formal", "Use formal, professional language"},
		{"casual", "Use a casual, conversational tone"},
		{"technical", "Use precise technical terminology"},
	}

	for _, tc := range styleTests {
		t.Run("Apply_"+tc.style, func(t *testing.T) {
			var spec prompt.PromptSpec
			rule.Apply(&spec, slots.Slots{Style: tc.style}, tr)
			if len(spec.Constraints) != 1 || spec.Constraints[0] != tc.want {
				t.Errorf("Constraints = %v, want [%q]", spec.Constraints, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ConstraintsFromEntitiesRule
// ---------------------------------------------------------------------------

func TestConstraintsFromEntitiesRule(t *testing.T) {
	rule := ConstraintsFromEntitiesRule()
	tr := loadTestTranslator(t)

	t.Run("When_entities_present", func(t *testing.T) {
		s := slots.Slots{Entities: []slots.Entity{{Text: "Go", Role: "implementation_medium"}}}
		if !rule.When(s) {
			t.Error("expected true")
		}
	})

	t.Run("When_entities_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_implementation_medium", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{Entities: []slots.Entity{{Text: "React", Role: "implementation_medium"}}}
		rule.Apply(&spec, s, tr)
		if len(spec.Constraints) != 1 || spec.Constraints[0] != "Focus on practical usage of React" {
			t.Errorf("Constraints = %v", spec.Constraints)
		}
	})

	t.Run("Apply_non_implementation_medium", func(t *testing.T) {
		var spec prompt.PromptSpec
		s := slots.Slots{Entities: []slots.Entity{{Text: "users", Role: "target_object"}}}
		rule.Apply(&spec, s, tr)
		if len(spec.Constraints) != 0 {
			t.Errorf("Constraints = %v, want empty", spec.Constraints)
		}
	})
}

// ---------------------------------------------------------------------------
// OutputFromFormatRule
// ---------------------------------------------------------------------------

func TestOutputFromFormatRule(t *testing.T) {
	rule := OutputFromFormatRule()
	tr := loadTestTranslator(t)

	t.Run("When_format_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Format: "bullets"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_format_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	formatTests := []struct {
		format string
		want   string
	}{
		{"bullets", "Use bullet points"},
		{"steps", "Use a step-by-step structure"},
		{"prose", "Use flowing prose paragraphs"},
		{"table", "Use a table format for comparison"},
		{"code", "Include code examples with comments"},
	}

	for _, tc := range formatTests {
		t.Run("Apply_"+tc.format, func(t *testing.T) {
			var spec prompt.PromptSpec
			rule.Apply(&spec, slots.Slots{Format: tc.format}, tr)
			if len(spec.OutputSpec) != 1 || spec.OutputSpec[0] != tc.want {
				t.Errorf("OutputSpec = %v, want [%q]", spec.OutputSpec, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// OutputFromIntentRule
// ---------------------------------------------------------------------------

func TestOutputFromIntentRule(t *testing.T) {
	rule := OutputFromIntentRule()
	tr := loadTestTranslator(t)

	t.Run("When_intent_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "howto"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_intent_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	intentTests := []struct {
		intent string
		want   string
	}{
		{"howto", "Present the solution in a clear, sequential order"},
		{"explain", "Use clear section headings"},
		{"generate", "Present the output in a ready-to-use format"},
		{"analyze", "Structure findings with clear categories"},
		{"decide", "Present options side by side for easy comparison"},
	}

	for _, tc := range intentTests {
		t.Run("Apply_"+tc.intent, func(t *testing.T) {
			var spec prompt.PromptSpec
			rule.Apply(&spec, slots.Slots{Intent: tc.intent}, tr)
			if len(spec.OutputSpec) != 1 || spec.OutputSpec[0] != tc.want {
				t.Errorf("OutputSpec = %v, want [%q]", spec.OutputSpec, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// OutputFromDepthRule
// ---------------------------------------------------------------------------

func TestOutputFromDepthRule(t *testing.T) {
	rule := OutputFromDepthRule()
	tr := loadTestTranslator(t)

	t.Run("When_depth_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Depth: "short"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_depth_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_short", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "short"}, tr)
		if len(spec.OutputSpec) != 1 || spec.OutputSpec[0] != "Keep sections brief" {
			t.Errorf("OutputSpec = %v", spec.OutputSpec)
		}
	})

	t.Run("Apply_deep_no_output", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "deep"}, tr)
		if len(spec.OutputSpec) != 0 {
			t.Errorf("OutputSpec = %v, want empty", spec.OutputSpec)
		}
	})

	t.Run("Apply_standard_no_output", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Depth: "standard"}, tr)
		if len(spec.OutputSpec) != 0 {
			t.Errorf("OutputSpec = %v, want empty", spec.OutputSpec)
		}
	})
}

// ---------------------------------------------------------------------------
// OutputFallbackRule
// ---------------------------------------------------------------------------

func TestOutputFallbackRule(t *testing.T) {
	rule := OutputFallbackRule()
	tr := loadTestTranslator(t)

	t.Run("When_always_true", func(t *testing.T) {
		if !rule.When(slots.Slots{}) {
			t.Error("expected true always")
		}
	})

	t.Run("Apply_output_empty_sets_fallback", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{}, tr)
		if len(spec.OutputSpec) != 1 || spec.OutputSpec[0] != "Use a clear and readable structure" {
			t.Errorf("OutputSpec = %v", spec.OutputSpec)
		}
	})

	t.Run("Apply_output_not_empty_no_change", func(t *testing.T) {
		spec := prompt.PromptSpec{OutputSpec: []string{"existing"}}
		rule.Apply(&spec, slots.Slots{}, tr)
		if len(spec.OutputSpec) != 1 || spec.OutputSpec[0] != "existing" {
			t.Errorf("OutputSpec = %v", spec.OutputSpec)
		}
	})
}

// ---------------------------------------------------------------------------
// QualityBaseRule
// ---------------------------------------------------------------------------

func TestQualityBaseRule(t *testing.T) {
	rule := QualityBaseRule()
	tr := loadTestTranslator(t)

	t.Run("When_always_true", func(t *testing.T) {
		if !rule.When(slots.Slots{}) {
			t.Error("expected true always")
		}
	})

	t.Run("Apply", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{}, tr)
		if len(spec.QualityCriteria) != 2 {
			t.Fatalf("len(QualityCriteria) = %d, want 2", len(spec.QualityCriteria))
		}
		if spec.QualityCriteria[0] != "Clear and structured" {
			t.Errorf("QualityCriteria[0] = %q", spec.QualityCriteria[0])
		}
		if spec.QualityCriteria[1] != "Accurate" {
			t.Errorf("QualityCriteria[1] = %q", spec.QualityCriteria[1])
		}
	})
}

// ---------------------------------------------------------------------------
// QualityFromIntentRule
// ---------------------------------------------------------------------------

func TestQualityFromIntentRule(t *testing.T) {
	rule := QualityFromIntentRule()
	tr := loadTestTranslator(t)

	t.Run("When_intent_present", func(t *testing.T) {
		if !rule.When(slots.Slots{Intent: "generate"}) {
			t.Error("expected true")
		}
	})

	t.Run("When_intent_empty", func(t *testing.T) {
		if rule.When(slots.Slots{}) {
			t.Error("expected false")
		}
	})

	t.Run("Apply_generate", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "generate"}, tr)
		if len(spec.QualityCriteria) != 1 || spec.QualityCriteria[0] != "Complete and ready to use" {
			t.Errorf("QualityCriteria = %v", spec.QualityCriteria)
		}
	})

	t.Run("Apply_analyze", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "analyze"}, tr)
		if len(spec.QualityCriteria) != 1 || spec.QualityCriteria[0] != "Balanced and evidence-based" {
			t.Errorf("QualityCriteria = %v", spec.QualityCriteria)
		}
	})

	t.Run("Apply_decide", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "decide"}, tr)
		if len(spec.QualityCriteria) != 1 || spec.QualityCriteria[0] != "Fair comparison of all options" {
			t.Errorf("QualityCriteria = %v", spec.QualityCriteria)
		}
	})

	t.Run("Apply_explain_nothing_added", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "explain"}, tr)
		if len(spec.QualityCriteria) != 0 {
			t.Errorf("QualityCriteria = %v, want empty", spec.QualityCriteria)
		}
	})

	t.Run("Apply_howto_nothing_added", func(t *testing.T) {
		var spec prompt.PromptSpec
		rule.Apply(&spec, slots.Slots{Intent: "howto"}, tr)
		if len(spec.QualityCriteria) != 0 {
			t.Errorf("QualityCriteria = %v, want empty", spec.QualityCriteria)
		}
	})
}
