package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Use project root as base dir
	baseDir := findProjectRoot(t)

	cfg, err := Load(baseDir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Verify intents loaded
	if len(cfg.Intents.Intents) != 8 {
		t.Errorf("expected 8 intents, got %d", len(cfg.Intents.Intents))
	}
	for _, name := range []string{"explain", "howto", "generate", "analyze", "decide", "debug", "refactor", "summarize"} {
		if _, ok := cfg.Intents.Intents[name]; !ok {
			t.Errorf("missing intent %q", name)
		}
	}

	// Verify explain has EN keywords
	explain := cfg.Intents.Intents["explain"]
	if len(explain.Keywords["en"]) == 0 {
		t.Error("explain intent has no EN keywords")
	}
	if len(explain.Phrases["en"]) == 0 {
		t.Error("explain intent has no EN phrases")
	}

	// Verify modifiers loaded
	if len(cfg.Modifiers.Audience) == 0 {
		t.Error("no audience modifiers loaded")
	}
	if len(cfg.Modifiers.Depth) == 0 {
		t.Error("no depth modifiers loaded")
	}
	if len(cfg.Modifiers.Style) == 0 {
		t.Error("no style modifiers loaded")
	}
	if len(cfg.Modifiers.Format) == 0 {
		t.Error("no format modifiers loaded")
	}

	// Verify stages loaded
	if len(cfg.Stages) != 3 {
		t.Errorf("expected 3 stages, got %d", len(cfg.Stages))
	}

	// Verify entities loaded
	if len(cfg.Entities) != 4 {
		t.Errorf("expected 4 entity roles, got %d", len(cfg.Entities))
	}

	// Verify languages loaded
	if len(cfg.Languages) < 2 {
		t.Errorf("expected at least 2 languages, got %d", len(cfg.Languages))
	}
	if en, ok := cfg.Languages["en"]; !ok {
		t.Error("missing English language config")
	} else if en.Name != "English" {
		t.Errorf("English name = %q, want %q", en.Name, "English")
	}
	if de, ok := cfg.Languages["de"]; !ok {
		t.Error("missing German language config")
	} else if de.Name != "German" {
		t.Errorf("German name = %q, want %q", de.Name, "German")
	}
}

func TestLoadMissingDir(t *testing.T) {
	_, err := Load("/nonexistent/path")
	if err == nil {
		t.Error("expected error for missing directory")
	}
}

func TestLoadMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	langDir := filepath.Join(dir, "languages")
	_ = os.MkdirAll(dataDir, 0o755)
	_ = os.MkdirAll(langDir, 0o755)

	// Write malformed intents file
	_ = os.WriteFile(filepath.Join(dataDir, "intents.yaml"), []byte("intents:\n  - bad: [structure\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "modifiers.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "stages.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "entities.yaml"), []byte("{}"), 0o644)

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error for malformed YAML")
	}
}

func TestFindBaseDirWithEnv(t *testing.T) {
	baseDir := findProjectRoot(t)
	t.Setenv("PROMPTC_DATA", baseDir)

	dir, err := FindBaseDir()
	if err != nil {
		t.Fatalf("FindBaseDir() error: %v", err)
	}
	if dir != baseDir {
		t.Errorf("FindBaseDir() = %q, want %q", dir, baseDir)
	}
}

func TestFindBaseDirBadEnv(t *testing.T) {
	t.Setenv("PROMPTC_DATA", "/nonexistent/path")
	_, err := FindBaseDir()
	if err == nil {
		t.Error("expected error for bad PROMPTC_DATA")
	}
}

func TestLoadAcronyms(t *testing.T) {
	baseDir := findProjectRoot(t)
	cfg, err := Load(baseDir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(cfg.Acronyms.Acronyms) == 0 {
		t.Error("expected acronyms to be loaded")
	}
	found := false
	for _, a := range cfg.Acronyms.Acronyms {
		if a == "REST" {
			found = true
		}
	}
	if !found {
		t.Error("expected REST in acronyms")
	}
}

func TestLoadMissingAcronyms(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	langDir := filepath.Join(dir, "languages")
	_ = os.MkdirAll(dataDir, 0o755)
	_ = os.MkdirAll(langDir, 0o755)

	_ = os.WriteFile(filepath.Join(dataDir, "intents.yaml"), []byte("intents: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "modifiers.yaml"), []byte("audience: {}\ndepth: {}\nstyle: {}\nformat: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "stages.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "entities.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(langDir, "en.yaml"), []byte("name: English\ncode: en\nstop_words: []\ntopic_verbs: []"), 0o644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() should succeed without acronyms.yaml: %v", err)
	}
	if len(cfg.Acronyms.Acronyms) != 0 {
		t.Errorf("expected empty acronyms, got %d", len(cfg.Acronyms.Acronyms))
	}
}

func TestLoadLanguageStopWords(t *testing.T) {
	baseDir := findProjectRoot(t)
	cfg, err := Load(baseDir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	en := cfg.Languages["en"]
	if len(en.StopWords) == 0 {
		t.Error("expected English stop words")
	}
	de := cfg.Languages["de"]
	if len(de.StopWords) == 0 {
		t.Error("expected German stop words")
	}
}

func TestFindBaseDir_FallbackToWorkingDir(t *testing.T) {
	t.Setenv("PROMPTC_DATA", "")
	// When run from project root, should find data/ dir
	dir, err := FindBaseDir()
	if err != nil {
		t.Skipf("not running from project root: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty directory")
	}
}

func TestIsDataDir_NonexistentPath(t *testing.T) {
	if isDataDir("/nonexistent/path/that/doesnt/exist") {
		t.Error("expected false for nonexistent path")
	}
}

func TestIsDataDir_FileNotDir(t *testing.T) {
	dir := t.TempDir()
	// Create a file named "data" (not a directory)
	_ = os.WriteFile(filepath.Join(dir, "data"), []byte("not a dir"), 0o644)
	if isDataDir(dir) {
		t.Error("expected false when data is a file, not a directory")
	}
}

func TestLoadYAML_NonexistentFile(t *testing.T) {
	var target map[string]string
	err := loadYAML("/nonexistent/file.yaml", &target)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadYAML_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	_ = os.WriteFile(path, []byte("{{invalid yaml"), 0o644)

	var target map[string]string
	err := loadYAML(path, &target)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_MissingLanguagesDir(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	_ = os.MkdirAll(dataDir, 0o755)
	// No languages/ dir

	_ = os.WriteFile(filepath.Join(dataDir, "intents.yaml"), []byte("intents: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "modifiers.yaml"), []byte("audience: {}\ndepth: {}\nstyle: {}\nformat: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "stages.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "entities.yaml"), []byte("{}"), 0o644)

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error when languages/ dir is missing")
	}
}

func TestLoad_LanguageDirSkipsNonYAML(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	langDir := filepath.Join(dir, "languages")
	_ = os.MkdirAll(dataDir, 0o755)
	_ = os.MkdirAll(langDir, 0o755)

	_ = os.WriteFile(filepath.Join(dataDir, "intents.yaml"), []byte("intents: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "modifiers.yaml"), []byte("audience: {}\ndepth: {}\nstyle: {}\nformat: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "stages.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "entities.yaml"), []byte("{}"), 0o644)
	// Write a non-yaml file and a subdirectory — both should be skipped
	_ = os.WriteFile(filepath.Join(langDir, "readme.txt"), []byte("not yaml"), 0o644)
	_ = os.MkdirAll(filepath.Join(langDir, "subdir"), 0o755)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(cfg.Languages) != 0 {
		t.Errorf("expected 0 languages, got %d", len(cfg.Languages))
	}
}

func TestLoad_LanguageWithoutCode(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	langDir := filepath.Join(dir, "languages")
	_ = os.MkdirAll(dataDir, 0o755)
	_ = os.MkdirAll(langDir, 0o755)

	_ = os.WriteFile(filepath.Join(dataDir, "intents.yaml"), []byte("intents: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "modifiers.yaml"), []byte("audience: {}\ndepth: {}\nstyle: {}\nformat: {}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "stages.yaml"), []byte("{}"), 0o644)
	_ = os.WriteFile(filepath.Join(dataDir, "entities.yaml"), []byte("{}"), 0o644)
	// Language file without code field — should derive from filename
	_ = os.WriteFile(filepath.Join(langDir, "fr.yaml"), []byte("name: French\nstop_words: []\ntopic_verbs: []"), 0o644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if _, ok := cfg.Languages["fr"]; !ok {
		t.Error("expected 'fr' language derived from filename")
	}
}

func TestLoadPhrases(t *testing.T) {
	baseDir := findProjectRoot(t)
	cfg, err := Load(baseDir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(cfg.Phrases.Phrases) == 0 {
		t.Error("expected phrases to be loaded")
	}
	if len(cfg.Phrases.Phrases["en"]) == 0 {
		t.Error("expected English phrases")
	}
	if len(cfg.Phrases.Phrases["de"]) == 0 {
		t.Error("expected German phrases")
	}
	found := false
	for _, p := range cfg.Phrases.Phrases["en"] {
		if p == "machine learning" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'machine learning' in English phrases")
	}
	found = false
	for _, p := range cfg.Phrases.Phrases["de"] {
		if p == "maschinelles lernen" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'maschinelles lernen' in German phrases")
	}
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from test file location to find project root (has go.mod)
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
