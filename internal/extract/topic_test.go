package extract

import (
	"testing"
)

func TestExtractTopic(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{"explain single", []string{"explain", "closures"}, "closures"},
		{"explain with audience", []string{"explain", "closures", "for", "beginners"}, "closures"},
		{"describe with entity", []string{"describe", "the", "api", "with", "python"}, "the api"},
		{"generate with entity", []string{"generate", "a", "rest", "api", "with", "python"}, "a rest api"},
		{"create with entity", []string{"create", "a", "function", "using", "go"}, "a function"},
		{"analyze stops at modifier", []string{"analyze", "this", "code"}, "this"}, // "code" is a format modifier keyword
		{"howto phrase", []string{"how", "do", "i", "test", "with", "jest"}, "test"},
		{"howto start project", []string{"how", "do", "i", "start", "a", "project", "with", "php"}, "start a project"},
		{"should i use phrase", []string{"should", "i", "use", "react", "or", "vue"}, "react or vue"},
		{"fallback keeps stop words", []string{"the", "goal", "is", "clarity"}, "the goal is clarity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.ExtractTopic(tt.tokens)
			if got != tt.want {
				t.Errorf("ExtractTopic(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}
