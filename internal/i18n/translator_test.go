package i18n

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestLoadAndGet(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
  howto: "How to %s."
scope:
  explain_main: "Explain the main concept"
`)

	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if got := tr.Get("objective.explain"); got != "Explain %s." {
		t.Errorf("Get(objective.explain) = %q, want %q", got, "Explain %s.")
	}
	if got := tr.Get("scope.explain_main"); got != "Explain the main concept" {
		t.Errorf("Get(scope.explain_main) = %q, want %q", got, "Explain the main concept")
	}
}

func TestGetf(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
`)

	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if got := tr.Getf("objective.explain", "closures"); got != "Explain closures." {
		t.Errorf("Getf() = %q, want %q", got, "Explain closures.")
	}
}

func TestFallbackToEnglish(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
  howto: "How to %s."
`)
	writeTestYAML(t, dir, "de.yaml", `
objective:
  explain: "Erkläre %s."
`)

	tr, err := Load(dir, "de", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Present in DE
	if got := tr.Get("objective.explain"); got != "Erkläre %s." {
		t.Errorf("Get(objective.explain) = %q, want DE", got)
	}
	// Missing in DE, falls back to EN
	if got := tr.Get("objective.howto"); got != "How to %s." {
		t.Errorf("Get(objective.howto) = %q, want EN fallback", got)
	}
}

func TestMissingKey(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
`)

	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if got := tr.Get("nonexistent.key"); got != "nonexistent.key" {
		t.Errorf("Get(missing) = %q, want key echoed back", got)
	}
}

func TestLang(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
`)

	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if got := tr.Lang(); got != "en" {
		t.Errorf("Lang() = %q, want %q", got, "en")
	}
}

func TestLoadProjectTranslations(t *testing.T) {
	dir := findProjectRoot(t)
	trDir := filepath.Join(dir, "translations")

	// Load English
	en, err := Load(trDir, "en", "en")
	if err != nil {
		t.Fatalf("Load EN: %v", err)
	}

	// Spot-check keys exist
	requiredKeys := []string{
		"labels.objective", "labels.scope", "labels.constraints",
		"objective.explain", "objective.generate", "objective.decide",
		"scope.explain_main", "scope.generate_complete",
		"constraints.concise", "constraints.formal",
		"output.bullets", "output.howto",
		"quality.clear", "quality.accurate",
		"role.from_entity",
		"context.getting_started",
	}
	for _, key := range requiredKeys {
		if got := en.Get(key); got == key {
			t.Errorf("EN missing key %q", key)
		}
	}

	// Load German with EN fallback
	de, err := Load(trDir, "de", "en")
	if err != nil {
		t.Fatalf("Load DE: %v", err)
	}

	// Verify DE has its own translations (not just falling back)
	if got := de.Get("objective.explain"); got == en.Get("objective.explain") {
		t.Errorf("DE objective.explain should differ from EN")
	}
}

func TestLoadMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `bad: [yaml: structure`)
	_, err := Load(dir, "en", "en")
	if err == nil {
		t.Error("expected error for malformed YAML")
	}
}

func TestLoadMissingFallbackFile(t *testing.T) {
	dir := t.TempDir()
	// No files at all
	_, err := Load(dir, "en", "en")
	if err == nil {
		t.Error("expected error when fallback file missing")
	}
}

func TestLoadMissingTargetLanguage(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
`)
	// No fr.yaml exists
	tr, err := Load(dir, "fr", "en")
	if err != nil {
		t.Fatalf("should not error when target missing: %v", err)
	}
	// Should fall back to English
	if got := tr.Get("objective.explain"); got != "Explain %s." {
		t.Errorf("expected EN fallback, got %q", got)
	}
	if got := tr.Lang(); got != "fr" {
		t.Errorf("Lang() = %q, want %q", got, "fr")
	}
}

func TestLoadEmptyYAML(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", "")
	tr, err := Load(dir, "en", "en")
	if err != nil {
		// Empty YAML may cause a parse error with yaml.Unmarshal into map type
		t.Logf("Load returned error for empty file (acceptable): %v", err)
		return
	}
	// If it succeeds, missing keys should echo back
	if got := tr.Get("any.key"); got != "any.key" {
		t.Errorf("expected key echoed back for empty file, got %q", got)
	}
}

func TestGetfNoArgs(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
quality:
  clear: "Clear and structured"
`)
	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got := tr.Getf("quality.clear"); got != "Clear and structured" {
		t.Errorf("Getf no args = %q, want %q", got, "Clear and structured")
	}
}

func TestGetfMissingKeyWithArgs(t *testing.T) {
	dir := t.TempDir()
	writeTestYAML(t, dir, "en.yaml", `
objective:
  explain: "Explain %s."
`)
	tr, err := Load(dir, "en", "en")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	// Missing key echoed back, then Sprintf applied to it
	got := tr.Getf("missing.key", "arg")
	if got == "" {
		t.Error("Getf with missing key returned empty string")
	}
}

func TestTranslationSymmetry(t *testing.T) {
	dir := findProjectRoot(t)
	translationsDir := filepath.Join(dir, "translations")

	en, err := loadTranslationFile(filepath.Join(translationsDir, "en.yaml"))
	if err != nil {
		t.Fatalf("loading en.yaml: %v", err)
	}
	de, err := loadTranslationFile(filepath.Join(translationsDir, "de.yaml"))
	if err != nil {
		t.Fatalf("loading de.yaml: %v", err)
	}

	var missingInDe []string
	for key := range en {
		if _, ok := de[key]; !ok {
			missingInDe = append(missingInDe, key)
		}
	}
	sort.Strings(missingInDe)
	for _, key := range missingInDe {
		t.Errorf("key %q exists in en.yaml but missing in de.yaml", key)
	}

	var missingInEn []string
	for key := range de {
		if _, ok := en[key]; !ok {
			missingInEn = append(missingInEn, key)
		}
	}
	sort.Strings(missingInEn)
	for _, key := range missingInEn {
		t.Errorf("key %q exists in de.yaml but missing in en.yaml", key)
	}
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

func writeTestYAML(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}
