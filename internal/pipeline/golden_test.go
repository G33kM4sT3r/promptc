package pipeline

import (
	"os"
	"path/filepath"
	"promptc/internal/render"
	"strings"
	"testing"
)

func TestGolden(t *testing.T) {
	cfg := loadTestConfig(t)

	tests := []struct {
		name   string
		input  string
		lang   string // expected language for renderer selection
		golden string // filename in testdata/
	}{
		{
			name:   "explain EN beginners",
			input:  "explain closures for beginners",
			lang:   "en",
			golden: "explain_en_beginners.golden",
		},
		{
			name:   "generate EN python",
			input:  "generate a REST API with Python",
			lang:   "en",
			golden: "generate_en_python.golden",
		},
		{
			name:   "decide EN react vue",
			input:  "should I use React or Vue",
			lang:   "en",
			golden: "decide_en_react_vue.golden",
		},
		{
			name:   "howto EN start php",
			input:  "how do I start a project with PHP",
			lang:   "en",
			golden: "howto_en_start_php.golden",
		},
		{
			name:   "analyze EN java",
			input:  "analyze authentication code with Java",
			lang:   "en",
			golden: "analyze_en_java.golden",
		},
		{
			name:   "explain DE deep",
			input:  "erkläre dependency injection detailliert",
			lang:   "de",
			golden: "explain_de_deep.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := loadTestTranslator(t, tt.lang)
			s, spec, err := RunWithRules(tt.input, cfg, translator, nil)
			if err != nil {
				t.Fatalf("RunWithRules() error: %v", err)
			}

			_ = s // language already handled by translator
			r := render.NewTranslated(translator)
			got := r.Render(spec)

			goldenPath := filepath.Join("testdata", tt.golden)
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("reading golden file %s: %v", goldenPath, err)
			}

			wantStr := strings.TrimSpace(string(want))
			gotStr := strings.TrimSpace(got)

			if gotStr != wantStr {
				t.Errorf("output mismatch for %q\n\ngot:\n%s\n\nwant:\n%s", tt.input, gotStr, wantStr)
			}
		})
	}
}
