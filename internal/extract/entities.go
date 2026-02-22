package extract

import (
	"promptc/internal/config"
	"promptc/internal/slots"
	"strings"
)

// ExtractEntities extracts entities from tokenized input using preposition-based detection.
// Supports all entity roles from config with multi-word capture.
func (e *Extractor) ExtractEntities(tokens []string) []slots.Entity {
	var entities []slots.Entity
	joined := strings.Join(tokens, " ")

	// Check multi-word prepositions first (e.g. "such as", "wie zum beispiel")
	for prep, role := range e.prepToRole {
		if !strings.Contains(prep, " ") {
			continue // handled below
		}
		idx := strings.Index(joined, prep+" ")
		if idx < 0 {
			continue
		}
		after := joined[idx+len(prep)+1:]
		text := e.captureEntityText(after, role)
		if text != "" {
			entities = append(entities, slots.Entity{Text: text, Role: role})
		}
	}

	// Check single-word prepositions
	for i := 0; i < len(tokens)-1; i++ {
		role, ok := e.prepToRole[tokens[i]]
		if !ok {
			continue
		}
		// Skip multi-word prepositions (already handled)
		if strings.Contains(tokens[i], " ") {
			continue
		}

		// Capture tokens after preposition until stop word or another preposition
		var parts []string
		entry := e.roleConfig[role]
		stopWords := e.collectStopWords(entry)

		for j := i + 1; j < len(tokens); j++ {
			tok := tokens[j]
			lower := strings.ToLower(tok)
			// Stop at another preposition
			if _, isPrep := e.prepToRole[lower]; isPrep {
				break
			}
			// Stop at stop words for this role
			if stopWords[lower] {
				break
			}
			// Stop at modifier keywords
			if e.allModifierKeywords[lower] {
				break
			}
			parts = append(parts, tok)
			if !entry.MultiWord {
				break
			}
		}

		if len(parts) > 0 {
			entities = append(entities, slots.Entity{
				Text: strings.Join(parts, " "),
				Role: role,
			})
		}
	}

	return entities
}

func (e *Extractor) captureEntityText(after, role string) string {
	tokens := strings.Fields(after)
	entry := e.roleConfig[role]
	stopWords := e.collectStopWords(entry)

	var parts []string
	for _, tok := range tokens {
		lower := strings.ToLower(tok)
		if _, isPrep := e.prepToRole[lower]; isPrep {
			break
		}
		if stopWords[lower] {
			break
		}
		if e.allModifierKeywords[lower] {
			break
		}
		parts = append(parts, tok)
		if !entry.MultiWord {
			break
		}
	}
	return strings.Join(parts, " ")
}

func (e *Extractor) collectStopWords(entry config.EntityEntry) map[string]bool {
	sw := make(map[string]bool)
	for _, words := range entry.StopWords {
		for _, w := range words {
			sw[w] = true
		}
	}
	return sw
}
