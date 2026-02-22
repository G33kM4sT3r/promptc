package extract

import (
	"testing"
)

func TestDetectAudience(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"explain", "for", "beginners"}, "beginners"},
		{[]string{"explain", "for", "beginner"}, "beginners"},
		{[]string{"explain", "for", "novice"}, "beginners"},
		{[]string{"explain", "advanced", "topics"}, "advanced"},
		{[]string{"explain", "for", "expert"}, "advanced"},
		{[]string{"anfänger", "erklärung"}, "beginners"},
		{[]string{"explain", "closures"}, ""},
	}

	for _, tt := range tests {
		got := ext.DetectAudience(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectAudience(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}

func TestDetectDepth(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"explain", "closures", "brief"}, "short"},
		{[]string{"quick", "overview"}, "short"},
		{[]string{"detailed", "explanation"}, "deep"},
		{[]string{"comprehensive", "guide"}, "deep"},
		{[]string{"in-depth", "analysis"}, "deep"},
		{[]string{"kurz", "erkläre"}, "short"},
		{[]string{"detailliert", "erkläre"}, "deep"},
		{[]string{"explain", "closures"}, ""},
	}

	for _, tt := range tests {
		got := ext.DetectDepth(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectDepth(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}

func TestDetectStyle(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"formal", "explanation"}, "formal"},
		{[]string{"academic", "paper"}, "formal"},
		{[]string{"casual", "guide"}, "casual"},
		{[]string{"technical", "documentation"}, "technical"},
		{[]string{"wissenschaftlich", "erkläre"}, "formal"},
		{[]string{"explain", "closures"}, ""},
	}

	for _, tt := range tests {
		got := ext.DetectStyle(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectStyle(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}

func TestDetectFormat(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"explain", "as", "bullets"}, "bullets"},
		{[]string{"explain", "as", "list"}, "bullets"},
		{[]string{"step-by-step", "guide"}, "steps"},
		{[]string{"show", "as", "table"}, "table"},
		{[]string{"write", "code"}, "code"},
		{[]string{"liste", "von"}, "bullets"},
		{[]string{"explain", "closures"}, ""},
	}

	for _, tt := range tests {
		got := ext.DetectFormat(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectFormat(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}
