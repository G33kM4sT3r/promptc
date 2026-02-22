package extract

import "strings"

// DetectIntent identifies the user's intent from tokenized input.
// Multi-word phrases are checked first (longest match wins), then single keywords.
func (e *Extractor) DetectIntent(tokens []string) string {
	joined := strings.Join(tokens, " ")

	// Check multi-word phrases first (sorted by length desc for greedy match)
	for _, pm := range e.intentPhrases {
		if strings.Contains(joined, pm.Phrase) {
			return pm.Value
		}
	}

	// Check single keywords
	for _, t := range tokens {
		if intent, ok := e.intentKeywords[t]; ok {
			return intent
		}
	}

	return "unknown"
}
