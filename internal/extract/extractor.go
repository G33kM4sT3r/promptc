package extract

import (
	"promptc/internal/config"
	"sort"
)

// phraseMapping maps a multi-word phrase to an intent/value.
type phraseMapping struct {
	Phrase string
	Value  string
}

// Extractor holds pre-built lookup maps for config-driven extraction.
type Extractor struct {
	cfg *config.Config

	// intent lookups
	intentKeywords map[string]string // keyword → intent name
	intentPhrases  []phraseMapping   // sorted by phrase length desc

	// stage lookups
	stageKeywords map[string]string // keyword → stage name

	// modifier lookups
	audienceKeywords map[string]string // keyword → audience value
	depthKeywords    map[string]string // keyword → depth value
	styleKeywords    map[string]string // keyword → style value
	formatKeywords   map[string]string // keyword → format value

	// entity lookups
	prepToRole map[string]string // preposition → entity role
	roleConfig map[string]config.EntityEntry

	// all intent keywords (for topic extraction stop detection)
	allIntentKeywords map[string]bool

	// all prepositions (for topic extraction stop detection)
	allPrepositions map[string]bool

	// all modifier keywords (for topic extraction stop detection)
	allModifierKeywords map[string]bool
}

// New creates an Extractor with pre-built lookup maps from the config.
func New(cfg *config.Config) *Extractor {
	e := &Extractor{
		cfg:                 cfg,
		intentKeywords:      make(map[string]string),
		stageKeywords:       make(map[string]string),
		audienceKeywords:    make(map[string]string),
		depthKeywords:       make(map[string]string),
		styleKeywords:       make(map[string]string),
		formatKeywords:      make(map[string]string),
		prepToRole:          make(map[string]string),
		roleConfig:          make(map[string]config.EntityEntry),
		allIntentKeywords:   make(map[string]bool),
		allPrepositions:     make(map[string]bool),
		allModifierKeywords: make(map[string]bool),
	}

	e.buildIntentMaps()
	e.buildStageMaps()
	e.buildModifierMaps()
	e.buildEntityMaps()

	return e
}

func (e *Extractor) buildIntentMaps() {
	for intentName, entry := range e.cfg.Intents.Intents {
		for _, keywords := range entry.Keywords {
			for _, kw := range keywords {
				e.intentKeywords[kw] = intentName
				e.allIntentKeywords[kw] = true
			}
		}
		for _, phrases := range entry.Phrases {
			for _, phrase := range phrases {
				e.intentPhrases = append(e.intentPhrases, phraseMapping{
					Phrase: phrase,
					Value:  intentName,
				})
			}
		}
	}
	// Sort phrases by length descending for greedy matching
	sort.Slice(e.intentPhrases, func(i, j int) bool {
		return len(e.intentPhrases[i].Phrase) > len(e.intentPhrases[j].Phrase)
	})
}

func (e *Extractor) buildStageMaps() {
	for stageName, entry := range e.cfg.Stages {
		for _, keywords := range entry.Keywords {
			for _, kw := range keywords {
				e.stageKeywords[kw] = stageName
			}
		}
	}
}

func (e *Extractor) buildModifierMaps() {
	buildModMap := func(entries map[string]config.ModifierEntry) map[string]string {
		m := make(map[string]string)
		for value, entry := range entries {
			for _, keywords := range entry.Keywords {
				for _, kw := range keywords {
					m[kw] = value
					e.allModifierKeywords[kw] = true
				}
			}
		}
		return m
	}

	e.audienceKeywords = buildModMap(e.cfg.Modifiers.Audience)
	e.depthKeywords = buildModMap(e.cfg.Modifiers.Depth)
	e.styleKeywords = buildModMap(e.cfg.Modifiers.Style)
	e.formatKeywords = buildModMap(e.cfg.Modifiers.Format)
}

func (e *Extractor) buildEntityMaps() {
	for roleName, entry := range e.cfg.Entities {
		e.roleConfig[roleName] = entry
		for _, preps := range entry.Prepositions {
			for _, prep := range preps {
				// Multi-word prepositions stored as-is
				e.prepToRole[prep] = roleName
				e.allPrepositions[prep] = true
			}
		}
	}
}

// isStopToken returns true if the token is a preposition or modifier keyword
// (used during topic extraction to know when to stop collecting topic tokens).
func (e *Extractor) isStopToken(token string) bool {
	return e.allPrepositions[token] || e.allModifierKeywords[token]
}
