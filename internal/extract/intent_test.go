package extract

import (
	"promptc/internal/config"
	"testing"
)

func TestDetectIntent(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"explain", "closures"}, "explain"},
		{[]string{"describe", "the", "api"}, "explain"},
		{[]string{"define", "polymorphism"}, "explain"},
		{[]string{"how", "do", "i", "test"}, "howto"},
		{[]string{"generate", "a", "rest", "api"}, "generate"},
		{[]string{"create", "a", "function"}, "generate"},
		{[]string{"write", "a", "test"}, "generate"},
		{[]string{"build", "a", "server"}, "generate"},
		{[]string{"analyze", "this", "code"}, "analyze"},
		{[]string{"review", "the", "implementation"}, "analyze"},
		{[]string{"decide", "between", "options"}, "decide"},
		{[]string{"compare", "react", "vue"}, "decide"},
		{[]string{"choose", "a", "framework"}, "decide"},
		// German
		{[]string{"erkläre", "closures"}, "explain"},
		{[]string{"erstelle", "eine", "funktion"}, "generate"},
		{[]string{"analysiere", "den", "code"}, "analyze"},
		// Phrases
		{[]string{"what", "is", "a", "closure"}, "explain"},
		{[]string{"how", "do", "i", "start"}, "howto"},
		{[]string{"should", "i", "use", "react"}, "decide"},
		{[]string{"which", "is", "better", "react", "vue"}, "decide"},
		// Unknown
		{[]string{"hello"}, "unknown"},
		{[]string{}, "unknown"},
	}

	for _, tt := range tests {
		got := ext.DetectIntent(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectIntent(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	baseDir := findProjectRoot(t)
	cfg, err := config.Load(baseDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
