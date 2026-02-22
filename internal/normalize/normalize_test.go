package normalize

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  hello world  ", "hello world"},
		{"what is a closure?", "what is a closure"},
		{"explain Node.js", "explain Node.js"},
		{"do this!", "do this"},
		{"end.", "end"},
		{"test...", "test"},
		{"\u201chello\u201d", "\"hello\""},
		{"\u2018hi\u2019", "'hi'"},
		{"foo\u2014bar", "foo-bar"},
		{"foo\u2013bar", "foo-bar"},
	}

	for _, tt := range tests {
		got := Normalize(tt.input)
		if got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExpandContractions(t *testing.T) {
	contractions := map[string]string{
		"don't":  "do not",
		"it's":   "it is",
		"can't":  "cannot",
		"i'm":    "i am",
		"what's": "what is",
	}
	tests := []struct {
		input string
		want  string
	}{
		{"don't do this", "do not do this"},
		{"it's a closure", "it is a closure"},
		{"I can't explain", "I cannot explain"},
		{"what's a monad", "what is a monad"},
		{"no contractions here", "no contractions here"},
		{"DON'T shout", "do not shout"},
	}
	for _, tt := range tests {
		got := ExpandContractions(tt.input, contractions)
		if got != tt.want {
			t.Errorf("ExpandContractions(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeAll(t *testing.T) {
	contractions := map[string]string{
		"don't": "do not",
	}
	got := NormalizeAll("  don\u2019t do this?  ", contractions)
	// After smart quote normalization, "don't" matches and expands
	want := "do not do this"
	if got != want {
		t.Errorf("NormalizeAll = %q, want %q", got, want)
	}
}
