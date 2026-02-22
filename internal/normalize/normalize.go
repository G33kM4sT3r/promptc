package normalize

import "strings"

// Normalize cleans input text for downstream processing.
// Strips trailing sentence punctuation while preserving dots within words (e.g. "Node.js").
// Normalizes smart quotes and em-dashes.
func Normalize(input string) string {
	s := strings.TrimSpace(input)

	// Normalize smart quotes to ASCII
	s = strings.ReplaceAll(s, "\u201c", "\"") // left double quote
	s = strings.ReplaceAll(s, "\u201d", "\"") // right double quote
	s = strings.ReplaceAll(s, "\u2018", "'")  // left single quote
	s = strings.ReplaceAll(s, "\u2019", "'")  // right single quote

	// Normalize em-dash to regular dash
	s = strings.ReplaceAll(s, "\u2014", "-") // em-dash
	s = strings.ReplaceAll(s, "\u2013", "-") // en-dash

	// Strip trailing sentence punctuation
	s = strings.TrimRight(s, "?!.")

	return s
}

// ExpandContractions replaces contractions with their expanded forms.
// Matching is case-insensitive but the expansion preserves sentence flow.
func ExpandContractions(input string, contractions map[string]string) string {
	if len(contractions) == 0 {
		return input
	}
	words := strings.Fields(input)
	for i, word := range words {
		lower := strings.ToLower(word)
		if expanded, ok := contractions[lower]; ok {
			words[i] = expanded
		}
	}
	return strings.Join(words, " ")
}

// NormalizeAll applies full normalization: smart quotes, dashes, punctuation, and contraction expansion.
func NormalizeAll(input string, contractions map[string]string) string {
	s := Normalize(input)
	s = ExpandContractions(s, contractions)
	return s
}
