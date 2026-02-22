package language

import (
	"strings"

	"promptc/internal/config"

	fasttext "github.com/nano-interactive/go-fasttext"
)

const (
	detectionThreshold  = 2
	confidenceThreshold = 0.5
)

// Detector provides language detection using fastText or keyword fallback.
type Detector struct {
	model *fasttext.Model
	cfg   *config.Config
	useFT bool
}

// NewDetector creates a Detector. If modelPath points to a valid fastText model,
// uses ML-based detection. Otherwise falls back to keyword scoring.
func NewDetector(modelPath string, cfg *config.Config) (*Detector, error) {
	d := &Detector{cfg: cfg}

	if modelPath != "" {
		model, err := fasttext.Open(modelPath)
		if err == nil {
			d.model = &model
			d.useFT = true
		}
		// If open fails, silently fall back to heuristic
	}

	return d, nil
}

// Detect returns the ISO 639-1 language code or "unknown".
func (d *Detector) Detect(input string) string {
	if d.useFT {
		return d.detectFastText(input)
	}
	return d.detectKeywords(input)
}

// Close releases fastText model resources.
func (d *Detector) Close() error {
	if d.model != nil {
		return d.model.Close()
	}
	return nil
}

func (d *Detector) detectFastText(input string) string {
	pred, err := d.model.PredictOne(input, confidenceThreshold)
	if err != nil {
		return d.detectKeywords(input) // fallback on prediction error
	}
	// Labels are formatted as "__label__en"
	lang := strings.TrimPrefix(pred.Label, "__label__")
	return lang
}

func (d *Detector) detectKeywords(input string) string {
	tokens := strings.Fields(strings.ToLower(input))
	if len(tokens) == 0 {
		return "unknown"
	}

	tokenSet := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		tokenSet[t] = true
	}

	scores := make(map[string]int)

	// Score against language stop words
	for code, lang := range d.cfg.Languages {
		for _, sw := range lang.StopWords {
			if tokenSet[sw] {
				scores[code]++
			}
		}
	}

	// Score against intent keywords per language
	scoreKeywords := func(keywords map[string][]string) {
		for lang, kws := range keywords {
			for _, kw := range kws {
				if tokenSet[kw] {
					scores[lang]++
				}
			}
		}
	}

	for _, entry := range d.cfg.Intents.Intents {
		scoreKeywords(entry.Keywords)
	}

	// Score against modifier keywords per language
	scoreModifiers := func(entries map[string]config.ModifierEntry) {
		for _, entry := range entries {
			scoreKeywords(entry.Keywords)
		}
	}
	scoreModifiers(d.cfg.Modifiers.Audience)
	scoreModifiers(d.cfg.Modifiers.Depth)
	scoreModifiers(d.cfg.Modifiers.Style)
	scoreModifiers(d.cfg.Modifiers.Format)

	// Find highest scoring language
	bestCode := "unknown"
	bestScore := 0
	for code, score := range scores {
		if score > bestScore {
			bestScore = score
			bestCode = code
		}
	}

	if bestScore < detectionThreshold {
		return "unknown"
	}
	return bestCode
}

// Detect is a backward-compatible free function that uses keyword detection.
//
// Deprecated: Use NewDetector + Detector.Detect instead.
func Detect(input string, cfg *config.Config) string {
	d := &Detector{cfg: cfg}
	return d.detectKeywords(input)
}
