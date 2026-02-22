package main

import (
	"strings"
	"testing"
)

func TestCompileCommand_BasicInput(t *testing.T) {
	out, err := runBinary("compile", "explain kubernetes")
	if err != nil {
		t.Fatalf("compile failed: %v\n%s", err, out)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestCompileCommand_JSONOutput(t *testing.T) {
	out, err := runBinary("compile", "--output", "json", "explain docker")
	if err != nil {
		t.Fatalf("compile --output json failed: %v\n%s", err, out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON output starting with '{', got: %.50s", out)
	}
}

func TestCompileCommand_YAMLOutput(t *testing.T) {
	out, err := runBinary("compile", "--output", "yaml", "explain docker")
	if err != nil {
		t.Fatalf("compile --output yaml failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "role:") && !strings.Contains(out, "objective:") {
		t.Errorf("expected YAML output with role/objective, got: %.100s", out)
	}
}

func TestCompileCommand_ExplainMode(t *testing.T) {
	out, err := runBinary("compile", "--explain", "explain kubernetes")
	if err != nil {
		t.Fatalf("compile --explain failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Intent") && !strings.Contains(out, "intent") {
		t.Errorf("expected explain output with intent info, got: %.100s", out)
	}
}

func TestCompileCommand_ScoreMode(t *testing.T) {
	out, err := runBinary("compile", "--score", "explain kubernetes")
	if err != nil {
		t.Fatalf("compile --score failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Score") && !strings.Contains(out, "score") && !strings.Contains(out, "/100") {
		t.Errorf("expected score in output, got: %.100s", out)
	}
}

func TestCompileCommand_LangOverride(t *testing.T) {
	out, err := runBinary("compile", "--lang", "de", "erkläre Docker")
	if err != nil {
		t.Fatalf("compile --lang de failed: %v\n%s", err, out)
	}
	if out == "" {
		t.Error("expected non-empty output for German input")
	}
}

func TestCompileCommand_NoInput(t *testing.T) {
	_, err := runBinary("compile")
	if err == nil {
		t.Error("expected error for missing input")
	}
}
