package language

import (
	"os"
	"path/filepath"
	"promptc/internal/config"
	"testing"
)

func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func TestDetect(t *testing.T) {
	cfg := loadTestConfig(t)

	tests := []struct {
		input string
		want  string
	}{
		{"explain closures for beginners", "en"},
		{"how do I start a project with PHP", "en"},
		{"describe the API with Python", "en"},
		{"erkläre dependency injection detailliert", "de"},
		{"wie kann ich mit Go beginnen", "de"},
		{"erstelle eine Funktion mit Python", "de"},
		// Below threshold
		{"closures", "unknown"},
		{"", "unknown"},
	}

	for _, tt := range tests {
		got := Detect(tt.input, cfg)
		if got != tt.want {
			t.Errorf("Detect(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNewDetectorKeywordFallback(t *testing.T) {
	cfg := loadTestConfig(t)

	// Empty model path should create a keyword-only detector
	det, err := NewDetector("", cfg)
	if err != nil {
		t.Fatalf("NewDetector() error: %v", err)
	}
	defer func() { _ = det.Close() }()

	if det.useFT {
		t.Error("expected useFT=false for empty model path")
	}
}

func TestDetectorDetectKeywords(t *testing.T) {
	cfg := loadTestConfig(t)

	det, err := NewDetector("", cfg)
	if err != nil {
		t.Fatalf("NewDetector() error: %v", err)
	}
	defer func() { _ = det.Close() }()

	tests := []struct {
		input string
		want  string
	}{
		{"explain closures for beginners", "en"},
		{"erkläre dependency injection detailliert", "de"},
		{"closures", "unknown"},
		{"", "unknown"},
	}

	for _, tt := range tests {
		got := det.Detect(tt.input)
		if got != tt.want {
			t.Errorf("Detector.Detect(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNewDetectorInvalidModelPath(t *testing.T) {
	cfg := loadTestConfig(t)

	// Non-existent model path should fall back to keyword detection without error
	det, err := NewDetector("/nonexistent/model.ftz", cfg)
	if err != nil {
		t.Fatalf("NewDetector() error: %v", err)
	}
	defer func() { _ = det.Close() }()

	if det.useFT {
		t.Error("expected useFT=false for invalid model path")
	}

	// Should still work with keyword fallback
	got := det.Detect("explain closures for beginners")
	if got != "en" {
		t.Errorf("Detector.Detect() = %q, want %q", got, "en")
	}
}

func TestDetectorCloseNilModel(t *testing.T) {
	cfg := loadTestConfig(t)

	det, _ := NewDetector("", cfg)
	// Close should not panic with nil model
	if err := det.Close(); err != nil {
		t.Errorf("Close() error: %v", err)
	}
}
