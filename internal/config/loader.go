package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load reads all YAML configuration files from the given base directory.
// It expects data/ and languages/ subdirectories.
func Load(baseDir string) (*Config, error) {
	cfg := &Config{
		Languages: make(map[string]LanguageConfig),
	}

	// Load data files
	if err := loadYAML(filepath.Join(baseDir, "data", "intents.yaml"), &cfg.Intents); err != nil {
		return nil, fmt.Errorf("loading intents: %w", err)
	}
	if err := loadYAML(filepath.Join(baseDir, "data", "modifiers.yaml"), &cfg.Modifiers); err != nil {
		return nil, fmt.Errorf("loading modifiers: %w", err)
	}
	if err := loadYAML(filepath.Join(baseDir, "data", "stages.yaml"), &cfg.Stages); err != nil {
		return nil, fmt.Errorf("loading stages: %w", err)
	}
	if err := loadYAML(filepath.Join(baseDir, "data", "entities.yaml"), &cfg.Entities); err != nil {
		return nil, fmt.Errorf("loading entities: %w", err)
	}

	// Load acronyms (optional — don't error if missing)
	var acronyms AcronymsConfig
	acronymsPath := filepath.Join(baseDir, "data", "acronyms.yaml")
	if err := loadYAML(acronymsPath, &acronyms); err == nil {
		cfg.Acronyms = acronyms
	}

	// Load contractions (optional — don't error if missing)
	var contractions ContractionsConfig
	contractionsPath := filepath.Join(baseDir, "data", "contractions.yaml")
	if err := loadYAML(contractionsPath, &contractions); err == nil {
		cfg.Contractions = contractions
	}

	// Load phrases (optional — don't error if missing)
	var phrases PhrasesConfig
	phrasesPath := filepath.Join(baseDir, "data", "phrases.yaml")
	if err := loadYAML(phrasesPath, &phrases); err == nil {
		cfg.Phrases = phrases
	}

	// Load enrichments (optional — don't error if missing)
	var enrichments EnrichmentsConfig
	enrichmentsPath := filepath.Join(baseDir, "data", "enrichments.yaml")
	if err := loadYAML(enrichmentsPath, &enrichments); err == nil {
		cfg.Enrichments = enrichments
	}

	// Load language files
	langDir := filepath.Join(baseDir, "languages")
	entries, err := os.ReadDir(langDir)
	if err != nil {
		return nil, fmt.Errorf("reading languages directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}
		var lang LanguageConfig
		if err := loadYAML(filepath.Join(langDir, entry.Name()), &lang); err != nil {
			return nil, fmt.Errorf("loading language %s: %w", entry.Name(), err)
		}
		code := lang.Code
		if code == "" {
			// Derive code from filename (e.g. "en.yaml" → "en")
			code = entry.Name()[:len(entry.Name())-len(".yaml")]
		}
		cfg.Languages[code] = lang
	}

	return cfg, nil
}

// FindBaseDir resolves the data directory by checking (in order):
// 1. PROMPTC_DATA environment variable
// 2. Directory of the current executable
// 3. Current working directory
func FindBaseDir() (string, error) {
	// Check env var first
	if dir := os.Getenv("PROMPTC_DATA"); dir != "" {
		if isDataDir(dir) {
			return dir, nil
		}
		return "", fmt.Errorf("PROMPTC_DATA=%q does not contain data/ directory", dir)
	}

	// Check executable directory
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		if isDataDir(dir) {
			return dir, nil
		}
	}

	// Check working directory
	wd, err := os.Getwd()
	if err == nil {
		if isDataDir(wd) {
			return wd, nil
		}
	}

	return "", fmt.Errorf("cannot find data directory: set PROMPTC_DATA or run from project root")
}

func isDataDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "data"))
	return err == nil && info.IsDir()
}

func loadYAML(path string, target any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}
	return nil
}
