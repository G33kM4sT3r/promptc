package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"promptc/internal/config"
	"promptc/internal/i18n"
	"promptc/internal/slots"
)

func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()
	dir := findProjectRoot(t)
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg
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

func TestRunEmptyInput(t *testing.T) {
	cfg := loadTestConfig(t)
	_, err := Run("", cfg, nil)
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestRunExplainIntent(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("explain closures for beginners", cfg, nil)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if s.Intent != "explain" {
		t.Errorf("Intent = %q, want %q", s.Intent, "explain")
	}
	if s.Topic != "closures" {
		t.Errorf("Topic = %q, want %q", s.Topic, "closures")
	}
	if s.Audience != "beginners" {
		t.Errorf("Audience = %q, want %q", s.Audience, "beginners")
	}
}

func TestRunGenerateIntent(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("generate a REST API with Python", cfg, nil)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if s.Intent != "generate" {
		t.Errorf("Intent = %q, want %q", s.Intent, "generate")
	}
	if s.Topic == "" {
		t.Error("Topic should not be empty")
	}
	found := false
	for _, e := range s.Entities {
		if e.Role == "implementation_medium" && e.Text == "python" {
			found = true
		}
	}
	if !found {
		t.Error("expected implementation_medium entity 'python'")
	}
}

func loadTestTranslator(t *testing.T, lang string) *i18n.Translator {
	t.Helper()
	dir := findProjectRoot(t)
	if lang == "" || lang == "unknown" {
		lang = "en"
	}
	translator, err := i18n.Load(filepath.Join(dir, "languages"), lang, "en")
	if err != nil {
		t.Fatalf("failed to load translator: %v", err)
	}
	return translator
}

func TestRunWithRulesProducesOutput(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "en")
	_, spec, err := RunWithRules("explain closures for beginners", cfg, translator, nil)
	if err != nil {
		t.Fatalf("RunWithRules() error: %v", err)
	}
	if spec.Objective == "" {
		t.Error("Objective should not be empty")
	}
	if len(spec.Scope) == 0 {
		t.Error("Scope should not be empty")
	}
}

// --- Additional integration tests ---

func TestRunWithTrace(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "en")
	s, result, err := RunWithTrace("explain closures for beginners", cfg, translator, nil)
	if err != nil {
		t.Fatalf("RunWithTrace() error: %v", err)
	}
	if s.Intent != "explain" {
		t.Errorf("Intent = %q, want %q", s.Intent, "explain")
	}
	if len(result.Applied) == 0 {
		t.Error("expected some rules to fire")
	}
	if result.Spec.Objective == "" {
		t.Error("Objective should not be empty")
	}
}

func TestApplyRulesViaEngine(t *testing.T) {
	translator := loadTestTranslator(t, "en")
	s := slots.Slots{
		Intent:   "explain",
		Topic:    "closures",
		Audience: "beginners",
	}
	cfg := loadTestConfig(t)
	spec := ApplyRules(s, translator, cfg.Enrichments)
	if spec.Objective == "" {
		t.Error("Objective should not be empty")
	}
	if len(spec.Scope) == 0 {
		t.Error("Scope should not be empty")
	}
	if len(spec.Constraints) == 0 {
		t.Error("Constraints should not be empty for beginners")
	}
}

func TestApplyRulesWithTraceViaEngine(t *testing.T) {
	translator := loadTestTranslator(t, "en")
	cfg := loadTestConfig(t)
	s := slots.Slots{Intent: "generate", Topic: "REST API"}
	result := ApplyRulesWithTrace(s, translator, cfg.Enrichments)
	if result.Spec.Objective == "" {
		t.Error("Objective should not be empty")
	}
	if len(result.Applied) == 0 {
		t.Error("expected some rules to fire")
	}
	// objective.base should always fire when intent is set
	found := false
	for _, id := range result.Applied {
		if id == "objective.base" {
			found = true
		}
	}
	if !found {
		t.Error("expected objective.base to fire")
	}
}

func TestRunConstraintObject(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "en")
	s, _, err := RunWithRules("explain REST without frameworks", cfg, translator, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	found := false
	for _, e := range s.Entities {
		if e.Role == "constraint_object" && strings.Contains(e.Text, "frameworks") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected constraint_object entity about frameworks, got entities: %v", s.Entities)
	}
}

func TestRunTargetObject(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "en")
	s, _, err := RunWithRules("explain REST for web development", cfg, translator, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	found := false
	for _, e := range s.Entities {
		if e.Role == "target_object" && strings.Contains(e.Text, "web development") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected target_object entity about web development, got entities: %v", s.Entities)
	}
}

func TestRunGermanDetection(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "de")
	s, spec, err := RunWithRules("erkläre dependency injection detailliert", cfg, translator, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if s.Language != "de" {
		t.Errorf("Language = %q, want %q", s.Language, "de")
	}
	if spec.Objective == "" {
		t.Error("Objective should not be empty for German input")
	}
}

func TestRunTopicExtraction(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("generate a rest api", cfg, nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if s.Topic == "" {
		t.Error("Topic should not be empty")
	}
	if !strings.Contains(strings.ToLower(s.Topic), "rest api") {
		t.Errorf("Topic = %q, expected it to contain 'rest api'", s.Topic)
	}
}

func TestAllIntentsProduceOutput(t *testing.T) {
	cfg := loadTestConfig(t)
	translator := loadTestTranslator(t, "en")
	inputs := map[string]string{
		"explain":   "explain closures",
		"howto":     "how do I start a project",
		"generate":  "generate a REST API",
		"analyze":   "analyze this code",
		"decide":    "should I use React or Vue",
		"debug":     "debug why my API returns 500",
		"refactor":  "refactor this function for readability",
		"summarize": "summarize microservices architecture",
	}
	for intent, input := range inputs {
		t.Run(intent, func(t *testing.T) {
			_, spec, err := RunWithRules(input, cfg, translator, nil)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if spec.Objective == "" {
				t.Error("Objective should not be empty")
			}
			if len(spec.QualityCriteria) == 0 {
				t.Error("QualityCriteria should not be empty")
			}
		})
	}
}

func TestRunCalculatesTier(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("explain closures for beginners", cfg, nil)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if s.Tier == "" {
		t.Error("Tier should not be empty after Run()")
	}
	if s.Tier != "standard" {
		t.Errorf("Tier = %q, want %q (beginners = standard)", s.Tier, "standard")
	}
}

func TestRunTierDeep(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("explain closures in-depth", cfg, nil)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if s.Tier != "rich" {
		t.Errorf("Tier = %q, want %q (deep = rich)", s.Tier, "rich")
	}
}

func TestRunTierShort(t *testing.T) {
	cfg := loadTestConfig(t)
	s, err := Run("explain closures brief", cfg, nil)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if s.Tier != "minimal" {
		t.Errorf("Tier = %q, want %q (short = minimal)", s.Tier, "minimal")
	}
}

func TestRunWhitespaceOnlyInput(t *testing.T) {
	cfg := loadTestConfig(t)
	_, err := Run("   \t\n  ", cfg, nil)
	if err == nil {
		t.Error("expected error for whitespace-only input")
	}
}
