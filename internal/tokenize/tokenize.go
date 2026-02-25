package tokenize

import (
	"strconv"
	"strings"
)

// Tokenize splits input into lowercase tokens. Original simple fallback (Tier 3).
func Tokenize(input string) []string {
	return strings.Fields(strings.ToLower(input))
}

// TokenizeWithPhrases applies cascading tokenization:
// Tier 1: multi-word phrase matching, Tier 2: boundary-aware, Tier 3: simple split.
func TokenizeWithPhrases(input string, phrases []string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	lower := strings.ToLower(input)

	// Tier 1: Replace known phrases with placeholders, then tokenize
	type replacement struct {
		placeholder string
		phrase      string
	}
	var replacements []replacement

	// Sort phrases by length descending (longest match first)
	sorted := make([]string, len(phrases))
	copy(sorted, phrases)
	sortByLengthDesc(sorted)

	working := lower
	for i, phrase := range sorted {
		lp := strings.ToLower(phrase)
		if strings.Contains(working, lp) {
			ph := placeholderFor(i)
			working = strings.Replace(working, lp, ph, 1)
			replacements = append(replacements, replacement{placeholder: ph, phrase: lp})
		}
	}

	// Tier 2: boundary-aware tokenize the working string
	tokens := BoundaryTokenize(working)

	// Restore phrase placeholders
	for i, tok := range tokens {
		for _, r := range replacements {
			if tok == r.placeholder {
				tokens[i] = r.phrase
				break
			}
		}
	}

	return tokens
}

// BoundaryTokenize handles punctuation boundaries (Tier 2):
// - Strips possessives ('s, 's)
// - Preserves hyphenated compounds
// - Preserves quoted multi-word terms as single tokens
func BoundaryTokenize(input string) []string {
	lower := strings.ToLower(input)
	var tokens []string

	i := 0
	runes := []rune(lower)
	for i < len(runes) {
		// Skip whitespace
		if runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' {
			i++
			continue
		}

		// Handle quoted terms
		if runes[i] == '"' || runes[i] == '\u201c' || runes[i] == '\u201d' {
			end := findClosingQuote(runes, i)
			if end > i+1 {
				inner := strings.TrimSpace(string(runes[i+1 : end]))
				if inner != "" {
					tokens = append(tokens, inner)
				}
				i = end + 1
				continue
			}
			i++
			continue
		}

		// Collect word (letters, digits, hyphens)
		start := i
		for i < len(runes) && !isBreak(runes[i]) {
			i++
		}
		if i > start {
			word := string(runes[start:i])
			// Strip possessives
			word = stripPossessive(word)
			if word != "" {
				tokens = append(tokens, word)
			}
		} else {
			i++ // skip stray punctuation
		}
	}
	return tokens
}

func isBreak(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '"' || r == '\u201c' || r == '\u201d'
}

func findClosingQuote(runes []rune, start int) int {
	for i := start + 1; i < len(runes); i++ {
		if runes[i] == '"' || runes[i] == '\u201d' {
			return i
		}
	}
	return start // no closing quote found
}

func stripPossessive(word string) string {
	for _, suffix := range []string{"'s", "\u2019s"} {
		if strings.HasSuffix(word, suffix) {
			return word[:len(word)-len(suffix)]
		}
	}
	return word
}

func sortByLengthDesc(ss []string) {
	// Simple insertion sort — phrase lists are small
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && len(ss[j]) > len(ss[j-1]); j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}

func placeholderFor(i int) string {
	return "\x00ph" + strconv.Itoa(i) + "\x00"
}
