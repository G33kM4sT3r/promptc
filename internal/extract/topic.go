package extract

import "strings"

// ExtractTopic extracts the topic from tokenized input.
// It starts collecting after the intent phrase/keyword and stops at entity prepositions,
// modifier keywords, or stage keywords.
func (e *Extractor) ExtractTopic(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}

	// Determine where to start collecting (skip the intent trigger)
	startIdx := e.findTopicStart(tokens)
	if startIdx < 0 {
		// No intent found — use non-stop-word tokens as topic
		return e.fallbackTopic(tokens)
	}

	var topic []string
	for i := startIdx; i < len(tokens); i++ {
		t := tokens[i]
		// Stop at prepositions or modifier keywords
		if e.isStopToken(t) {
			break
		}
		topic = append(topic, t)
	}

	return strings.Join(topic, " ")
}

// findTopicStart returns the token index where topic collection should begin,
// or -1 if no intent trigger was found.
func (e *Extractor) findTopicStart(tokens []string) int {
	joined := strings.Join(tokens, " ")

	// Check phrases first (sorted by length desc)
	for _, pm := range e.intentPhrases {
		idx := strings.Index(joined, pm.Phrase)
		if idx >= 0 {
			// Find the token index after the phrase
			afterPhrase := joined[:idx] + strings.Repeat("x", len(pm.Phrase))
			_ = afterPhrase
			// Count how many tokens the phrase covers
			phraseTokenCount := len(strings.Fields(pm.Phrase))
			// Find which token index the phrase starts at
			prefixLen := 0
			startTokenIdx := 0
			for i, t := range tokens {
				tokStart := strings.Index(joined[prefixLen:], t) + prefixLen
				if tokStart >= idx {
					startTokenIdx = i
					break
				}
				prefixLen = tokStart + len(t)
			}
			return startTokenIdx + phraseTokenCount
		}
	}

	// Check single keywords
	for i, t := range tokens {
		if _, ok := e.intentKeywords[t]; ok {
			return i + 1
		}
	}

	return -1
}

func (e *Extractor) fallbackTopic(tokens []string) string {
	var topic []string
	for _, t := range tokens {
		if e.isStopToken(t) {
			break
		}
		topic = append(topic, t)
	}
	return strings.Join(topic, " ")
}

// CleanTopic strips leading articles (EN + DE) and applies acronym casing.
func (e *Extractor) CleanTopic(topic string) string {
	if topic == "" {
		return ""
	}

	words := strings.Fields(topic)
	if len(words) == 0 {
		return ""
	}

	// Strip leading articles (EN + DE)
	articles := map[string]bool{
		"a": true, "an": true, "the": true,
		"ein": true, "eine": true, "einen": true,
		"einer": true, "einem": true, "eines": true,
		"der": true, "die": true, "das": true,
	}
	for len(words) > 0 && articles[strings.ToLower(words[0])] {
		words = words[1:]
	}
	if len(words) == 0 {
		return ""
	}

	// Build acronym lookup
	acronymMap := make(map[string]string)
	for _, a := range e.cfg.Acronyms.Acronyms {
		acronymMap[strings.ToLower(a)] = a
	}

	// Apply acronym casing
	for i, w := range words {
		if canonical, ok := acronymMap[strings.ToLower(w)]; ok {
			words[i] = canonical
		}
	}

	return strings.Join(words, " ")
}
