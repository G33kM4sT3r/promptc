package extract

import (
	"testing"
)

func TestDetectStage(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		tokens []string
		want   string
	}{
		{[]string{"how", "to", "start", "a", "project"}, "getting-started"},
		{[]string{"how", "to", "begin", "with", "go"}, "getting-started"},
		{[]string{"implement", "authentication"}, "implementation"},
		{[]string{"optimize", "database", "queries"}, "optimization"},
		{[]string{"improve", "performance"}, "optimization"},
		// German
		{[]string{"beginnen", "mit", "go"}, "getting-started"},
		{[]string{"implementieren", "auth"}, "implementation"},
		{[]string{"optimieren", "datenbank"}, "optimization"},
		// No stage
		{[]string{"explain", "closures"}, ""},
	}

	for _, tt := range tests {
		got := ext.DetectStage(tt.tokens)
		if got != tt.want {
			t.Errorf("DetectStage(%v) = %q, want %q", tt.tokens, got, tt.want)
		}
	}
}
