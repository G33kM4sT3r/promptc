package render

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"promptc/internal/i18n"
	"promptc/internal/prompt"
)

// translationsDir returns the absolute path to the project's translations/ directory.
func translationsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "translations")
}

func TestTranslatedRendererEnglish(t *testing.T) {
	translator, err := i18n.Load(translationsDir(), "en", "en")
	if err != nil {
		t.Fatalf("failed to load English translations: %v", err)
	}

	r := NewTranslated(translator)

	spec := prompt.PromptSpec{
		Role:            "You are an expert in Go.",
		Objective:       "Explain closures.",
		Context:         "The explanation should be beginner-friendly.",
		Scope:           []string{"Explain the main concept"},
		Constraints:     []string{"Use simple language"},
		OutputSpec:      []string{"Use clear section headings"},
		QualityCriteria: []string{"Clear and structured"},
	}

	result := r.Render(spec)

	checks := []string{
		"You are an expert in Go.",
		"Objective:",
		"Explain closures.",
		"Context:",
		"beginner-friendly",
		"Scope:",
		"- Explain the main concept",
		"Constraints:",
		"- Use simple language",
		"Output:",
		"- Use clear section headings",
		"Quality criteria:",
		"- Clear and structured",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestTranslatedRendererGerman(t *testing.T) {
	translator, err := i18n.Load(translationsDir(), "de", "en")
	if err != nil {
		t.Fatalf("failed to load German translations: %v", err)
	}

	r := NewTranslated(translator)

	spec := prompt.PromptSpec{
		Role:            "Du bist ein Experte für Go.",
		Objective:       "Erkläre Closures.",
		Context:         "Die Erklärung sollte anfängerfreundlich sein.",
		Scope:           []string{"Das Hauptkonzept erklären"},
		Constraints:     []string{"Einfache Sprache verwenden"},
		OutputSpec:      []string{"Klare Überschriften verwenden"},
		QualityCriteria: []string{"Klar und strukturiert"},
	}

	result := r.Render(spec)

	checks := []string{
		"Du bist ein Experte für Go.",
		"Ziel:",
		"Erkläre Closures.",
		"Kontext:",
		"Umfang:",
		"Einschränkungen:",
		"Ausgabeformat:",
		"Qualitätskriterien:",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestTranslatedRendererEmptySpec(t *testing.T) {
	translator, err := i18n.Load(translationsDir(), "en", "en")
	if err != nil {
		t.Fatalf("failed to load translations: %v", err)
	}

	r := NewTranslated(translator)
	result := r.Render(prompt.PromptSpec{})

	if result != "" {
		t.Errorf("expected empty output for empty spec, got %q", result)
	}
}

func TestTranslatedRendererRenderScore(t *testing.T) {
	translator, err := i18n.Load(translationsDir(), "en", "en")
	if err != nil {
		t.Fatalf("failed to load translations: %v", err)
	}

	r := NewTranslated(translator)
	sb := ScoreBreakdown{
		Total:      65,
		Breakdown:  map[string]int{"objective": 25, "scope": 15, "output": 15, "quality": 10},
		MaxWeights: map[string]int{"role": 10, "objective": 25, "context": 15, "scope": 15, "constraints": 10, "output": 15, "quality": 10},
	}
	result := r.RenderScore(sb)

	for _, check := range []string{"Score: 65/100", "objective", "scope", "quality"} {
		if !strings.Contains(result, check) {
			t.Errorf("score output missing %q", check)
		}
	}
}
