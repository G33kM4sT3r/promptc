package extract

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Intent detection edge cases
// ---------------------------------------------------------------------------

func TestDetectIntent_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns unknown",
			tokens: []string{},
			want:   "unknown",
		},
		{
			name:   "nil tokens returns unknown",
			tokens: nil,
			want:   "unknown",
		},
		{
			name:   "single non-keyword token returns unknown",
			tokens: []string{"banana"},
			want:   "unknown",
		},
		{
			name:   "uppercase EXPLAIN not matched (tokens are lowercased by tokenizer)",
			tokens: []string{"EXPLAIN", "closures"},
			want:   "unknown",
		},
		{
			name:   "lowercase explain is matched",
			tokens: []string{"explain", "closures"},
			want:   "explain",
		},
		{
			name:   "phrase match takes priority over keyword position",
			tokens: []string{"explain", "how", "to", "generate", "code"},
			want:   "howto", // "how to" is a phrase match which is checked before single keywords
		},
		{
			name:   "multiple single keywords first token wins",
			tokens: []string{"explain", "closures", "analyze", "code"},
			want:   "explain",
		},
		{
			name:   "intent keyword not at start still detected",
			tokens: []string{"please", "explain", "this"},
			want:   "explain",
		},
		{
			name:   "single intent keyword alone",
			tokens: []string{"explain"},
			want:   "explain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectIntent(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectIntent(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Topic extraction edge cases
// ---------------------------------------------------------------------------

func TestExtractTopic_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "intent keyword at end with no topic after",
			tokens: []string{"explain"},
			want:   "",
		},
		{
			name:   "no intent keyword found uses fallback",
			tokens: []string{"banana", "smoothie"},
			want:   "banana smoothie",
		},
		{
			name:   "very long topic all captured until stop token",
			tokens: []string{"explain", "the", "very", "long", "topic", "covering", "many", "different", "things", "today", "and", "tomorrow"},
			want:   "the very long topic covering many different things today and tomorrow",
		},
		{
			name:   "topic with dots preserved",
			tokens: []string{"explain", "node.js", "middleware"},
			want:   "node.js middleware",
		},
		{
			name:   "topic with dots and entity after",
			tokens: []string{"explain", "node.js", "middleware", "with", "express"},
			want:   "node.js middleware",
		},
		{
			name:   "multiple tokens no intent fallback single token",
			tokens: []string{"redis"},
			want:   "redis",
		},
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

// ---------------------------------------------------------------------------
// Entity extraction edge cases
// ---------------------------------------------------------------------------

func TestExtractEntities_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name      string
		tokens    []string
		wantRoles map[string]string // role -> text
		wantCount int               // expected number of entities (-1 to skip check)
	}{
		{
			name:      "empty tokens no entities",
			tokens:    []string{},
			wantRoles: map[string]string{},
			wantCount: 0,
		},
		{
			name:      "nil tokens no entities",
			tokens:    nil,
			wantRoles: map[string]string{},
			wantCount: 0,
		},
		{
			name:   "multiple entities of different roles",
			tokens: []string{"build", "api", "with", "python", "for", "web", "development"},
			wantRoles: map[string]string{
				"implementation_medium": "python",
				"target_object":         "web development",
			},
			wantCount: 2,
		},
		{
			name:      "preposition at end of input no entity captured",
			tokens:    []string{"build", "api", "with"},
			wantRoles: map[string]string{},
			wantCount: 0,
		},
		{
			name:   "multi-word entity captured",
			tokens: []string{"build", "api", "with", "machine", "learning"},
			wantRoles: map[string]string{
				"implementation_medium": "machine learning",
			},
			wantCount: 1,
		},
		{
			name:      "no prepositions in input",
			tokens:    []string{"explain", "closures", "quickly"},
			wantRoles: map[string]string{},
			wantCount: 0,
		},
		{
			name:   "entity with constraint preposition",
			tokens: []string{"generate", "api", "without", "frameworks"},
			wantRoles: map[string]string{
				"constraint_object": "frameworks",
			},
			wantCount: 1,
		},
		{
			name:   "three different entity roles",
			tokens: []string{"build", "api", "with", "python", "for", "mobile", "without", "orm"},
			wantRoles: map[string]string{
				"implementation_medium": "python",
				"target_object":         "mobile",
				"constraint_object":     "orm",
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entities := ext.ExtractEntities(tt.tokens)
			got := make(map[string]string)
			for _, e := range entities {
				got[e.Role] = e.Text
			}

			if tt.wantCount >= 0 && len(entities) != tt.wantCount {
				t.Errorf("ExtractEntities(%v) returned %d entities, want %d; got %v",
					tt.tokens, len(entities), tt.wantCount, entities)
			}

			for role, wantText := range tt.wantRoles {
				gotText, ok := got[role]
				if !ok {
					t.Errorf("missing entity role %q in %v", role, entities)
				} else if gotText != wantText {
					t.Errorf("entity %q = %q, want %q", role, gotText, wantText)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Stage detection edge cases
// ---------------------------------------------------------------------------

func TestDetectStage_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "no stage keyword returns empty",
			tokens: []string{"explain", "closures"},
			want:   "",
		},
		{
			name:   "multiple stage keywords first one wins",
			tokens: []string{"start", "implement", "optimize"},
			want:   "getting-started",
		},
		{
			name:   "stage keyword not at beginning",
			tokens: []string{"please", "optimize", "the", "code"},
			want:   "optimization",
		},
		{
			name:   "single stage keyword alone",
			tokens: []string{"optimize"},
			want:   "optimization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectStage(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectStage(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Modifier detection edge cases
// ---------------------------------------------------------------------------

func TestDetectAudience_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "no audience keyword returns empty",
			tokens: []string{"explain", "closures"},
			want:   "",
		},
		{
			name:   "multiple audience keywords first one wins",
			tokens: []string{"beginners", "advanced", "expert"},
			want:   "beginners",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectAudience(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectAudience(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestDetectDepth_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "no depth keyword returns empty",
			tokens: []string{"explain", "closures"},
			want:   "",
		},
		{
			name:   "multiple depth keywords first one wins",
			tokens: []string{"brief", "detailed", "comprehensive"},
			want:   "short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectDepth(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectDepth(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestDetectStyle_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "no style keyword returns empty",
			tokens: []string{"explain", "closures"},
			want:   "",
		},
		{
			name:   "multiple style keywords first one wins",
			tokens: []string{"formal", "casual", "technical"},
			want:   "formal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectStyle(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectStyle(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestDetectFormat_EdgeCases(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "empty tokens returns empty",
			tokens: []string{},
			want:   "",
		},
		{
			name:   "nil tokens returns empty",
			tokens: nil,
			want:   "",
		},
		{
			name:   "no format keyword returns empty",
			tokens: []string{"explain", "closures"},
			want:   "",
		},
		{
			name:   "multiple format keywords first one wins",
			tokens: []string{"bullets", "step-by-step", "table"},
			want:   "bullets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ext.DetectFormat(tt.tokens)
			if got != tt.want {
				t.Errorf("DetectFormat(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Combined edge cases: multiple extractors on same input
// ---------------------------------------------------------------------------

func TestAllExtractors_EmptyInput(t *testing.T) {
	cfg := testConfig(t)
	ext := New(cfg)

	tokens := []string{}

	if got := ext.DetectIntent(tokens); got != "unknown" {
		t.Errorf("DetectIntent([]) = %q, want %q", got, "unknown")
	}
	if got := ext.ExtractTopic(tokens); got != "" {
		t.Errorf("ExtractTopic([]) = %q, want %q", got, "")
	}
	if got := ext.DetectStage(tokens); got != "" {
		t.Errorf("DetectStage([]) = %q, want %q", got, "")
	}
	if got := ext.DetectAudience(tokens); got != "" {
		t.Errorf("DetectAudience([]) = %q, want %q", got, "")
	}
	if got := ext.DetectDepth(tokens); got != "" {
		t.Errorf("DetectDepth([]) = %q, want %q", got, "")
	}
	if got := ext.DetectStyle(tokens); got != "" {
		t.Errorf("DetectStyle([]) = %q, want %q", got, "")
	}
	if got := ext.DetectFormat(tokens); got != "" {
		t.Errorf("DetectFormat([]) = %q, want %q", got, "")
	}
	if entities := ext.ExtractEntities(tokens); len(entities) != 0 {
		t.Errorf("ExtractEntities([]) = %v, want empty", entities)
	}
}
