package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Translator provides translated strings for a target language with English fallback.
type Translator struct {
	strings  map[string]string // flat "section.key" → translated string
	fallback map[string]string // English baseline
	lang     string
}

// Load creates a Translator for the given language, with fallback to fallbackLang.
// translationsDir contains YAML files named by language code (e.g., "en.yaml", "de.yaml").
func Load(translationsDir, lang, fallbackLang string) (*Translator, error) {
	fallbackStrings, err := loadTranslationFile(filepath.Join(translationsDir, fallbackLang+".yaml"))
	if err != nil {
		return nil, fmt.Errorf("loading fallback translations (%s): %w", fallbackLang, err)
	}

	var langStrings map[string]string
	if lang == fallbackLang {
		langStrings = fallbackStrings
	} else {
		langStrings, err = loadTranslationFile(filepath.Join(translationsDir, lang+".yaml"))
		if err != nil {
			// If target language file missing, use fallback entirely
			langStrings = make(map[string]string)
		}
	}

	return &Translator{
		strings:  langStrings,
		fallback: fallbackStrings,
		lang:     lang,
	}, nil
}

// Get returns the translated string for the given key.
// Falls back to English if not found in the target language.
// Returns the key itself if not found in any language.
func (t *Translator) Get(key string) string {
	if v, ok := t.strings[key]; ok {
		return v
	}
	if v, ok := t.fallback[key]; ok {
		return v
	}
	return key
}

// Getf returns the translated string with fmt.Sprintf formatting applied.
func (t *Translator) Getf(key string, args ...any) string {
	template := t.Get(key)
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}

// Lang returns the translator's target language code.
func (t *Translator) Lang() string {
	return t.lang
}

func loadTranslationFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse as nested map
	var raw map[string]map[string]string
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	// Flatten to "section.key" → value
	flat := make(map[string]string)
	for section, entries := range raw {
		for key, value := range entries {
			flat[section+"."+key] = value
		}
	}

	return flat, nil
}
