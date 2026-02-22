package extract

import (
	"testing"
)

func TestExtractEntities(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   map[string]string // role → text
	}{
		{
			name:   "with preposition",
			tokens: []string{"explain", "closures", "with", "python"},
			want:   map[string]string{"implementation_medium": "python"},
		},
		{
			name:   "using preposition",
			tokens: []string{"build", "api", "using", "go"},
			want:   map[string]string{"implementation_medium": "go"},
		},
		{
			name:   "for preposition",
			tokens: []string{"explain", "closures", "for", "web", "development"},
			want:   map[string]string{"target_object": "web development"},
		},
		{
			name:   "without preposition",
			tokens: []string{"generate", "code", "without", "frameworks"},
			want:   map[string]string{"constraint_object": "frameworks"},
		},
		{
			name:   "multi-word entity",
			tokens: []string{"build", "api", "with", "node.js", "express"},
			want:   map[string]string{"implementation_medium": "node.js express"},
		},
		{
			name:   "no entities",
			tokens: []string{"explain", "closures"},
			want:   map[string]string{},
		},
		{
			name:   "multi-word prep preserves original case",
			tokens: []string{"build", "api", "such", "as", "GraphQL"},
			want:   map[string]string{"example_object": "GraphQL"},
		},
		{
			name:   "single-word prep stops at capitalized stop word",
			tokens: []string{"build", "api", "with", "python", "And", "more"},
			want:   map[string]string{"implementation_medium": "python"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entities := ext.ExtractEntities(tt.tokens)
			got := make(map[string]string)
			for _, e := range entities {
				got[e.Role] = e.Text
			}
			for role, wantText := range tt.want {
				if gotText, ok := got[role]; !ok {
					t.Errorf("missing entity role %q", role)
				} else if gotText != wantText {
					t.Errorf("entity %q = %q, want %q", role, gotText, wantText)
				}
			}
		})
	}
}
