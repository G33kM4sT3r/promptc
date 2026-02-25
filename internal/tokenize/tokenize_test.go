package tokenize

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"EXPLAIN Closures", []string{"explain", "closures"}},
		{"  spaced  out  ", []string{"spaced", "out"}},
		{"Node.js API", []string{"node.js", "api"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		got := Tokenize(tt.input)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestTokenizeWithPhrases(t *testing.T) {
	phrases := []string{"machine learning", "REST API", "clean architecture"}

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		// Tier 1: phrase matching
		{"phrase match", "explain machine learning basics", []string{"explain", "machine learning", "basics"}},
		{"phrase case insensitive", "explain Machine Learning", []string{"explain", "machine learning"}},
		{"multiple phrases", "compare machine learning and clean architecture", []string{"compare", "machine learning", "and", "clean architecture"}},
		// Tier 2: boundary-aware
		{"possessive stripped", "explain python's gil", []string{"explain", "python", "gil"}},
		{"hyphenated compound", "explain well-known patterns", []string{"explain", "well-known", "patterns"}},
		{"quoted term", `explain "clean code" principles`, []string{"explain", "clean code", "principles"}},
		// Tier 3: fallback (no phrases, no special punctuation)
		{"simple fallback", "explain closures", []string{"explain", "closures"}},
		{"empty input", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TokenizeWithPhrases(tt.input, phrases)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenizeWithPhrases(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestBoundaryTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"python's gil", []string{"python", "gil"}},
		{"it's a test", []string{"it", "a", "test"}},
		{"well-known pattern", []string{"well-known", "pattern"}},
		{`use "clean code"`, []string{"use", "clean code"}},
		{`say "hello world" now`, []string{"say", "hello world", "now"}},
		{"simple words", []string{"simple", "words"}},
		{"trailing-dash-", []string{"trailing-dash-"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := BoundaryTokenize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoundaryTokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSortByLengthDesc(t *testing.T) {
	ss := []string{"ab", "abcde", "abc", "a"}
	sortByLengthDesc(ss)
	want := []string{"abcde", "abc", "ab", "a"}
	if !reflect.DeepEqual(ss, want) {
		t.Errorf("sortByLengthDesc = %v, want %v", ss, want)
	}
}
