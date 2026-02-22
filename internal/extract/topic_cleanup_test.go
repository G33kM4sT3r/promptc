package extract

import (
	"testing"

	"promptc/internal/config"
)

func TestCleanTopic(t *testing.T) {
	cfg := &config.Config{
		Acronyms: config.AcronymsConfig{
			Acronyms: []string{"REST", "API", "ORM", "PHP"},
		},
	}
	ext := New(cfg)

	tests := []struct {
		input string
		want  string
	}{
		{"a rest api", "REST API"},
		{"the database", "database"},
		{"an orm", "ORM"},
		{"closures", "closures"},
		{"", ""},
		{"a", ""},
		{"the", ""},
		{"php framework", "PHP framework"},
		{"ein rest api", "REST API"},
	}

	for _, tt := range tests {
		got := ext.CleanTopic(tt.input)
		if got != tt.want {
			t.Errorf("CleanTopic(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
