package explain

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"promptc/internal/rules"
	"promptc/internal/slots"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrint(t *testing.T) {
	s := slots.Slots{
		Language: "en",
		Intent:   "explain",
		Topic:    "REST APIs",
		Stage:    "getting-started",
		Entities: []slots.Entity{{Text: "Python", Role: "implementation_medium"}},
		Audience: "beginner",
		Depth:    "overview",
		Style:    "tutorial",
		Format:   "step-by-step",
	}
	trace := &rules.ApplyResult{
		Applied: []string{"objective.from_intent", "context.from_stage"},
		Skipped: []string{"scope.explain_topic"},
	}

	out := captureOutput(func() { Print(s, trace) })

	if !strings.Contains(out, "explain") {
		t.Error("expected intent in output")
	}
	if !strings.Contains(out, "REST APIs") {
		t.Error("expected topic in output")
	}
	if !strings.Contains(out, "objective.from_intent") {
		t.Error("expected applied rules in output")
	}
	if !strings.Contains(out, "scope.explain_topic") {
		t.Error("expected skipped rules in output")
	}
	if !strings.Contains(out, "Python") {
		t.Error("expected entity in output")
	}
}

func TestPrint_EmptySlots(t *testing.T) {
	s := slots.Slots{}

	out := captureOutput(func() { Print(s, nil) })

	if !strings.Contains(out, "<none>") {
		t.Error("expected <none> for empty values")
	}
}

func TestPrint_WithNilTrace(t *testing.T) {
	s := slots.Slots{Intent: "howto", Topic: "Docker"}

	out := captureOutput(func() { Print(s, nil) })

	if strings.Contains(out, "Rules:") {
		t.Error("expected no Rules section when trace is nil")
	}
}

func TestFormatRuleList_Empty(t *testing.T) {
	result := formatRuleList(nil)
	if result != "<none>" {
		t.Errorf("expected <none>, got %q", result)
	}
}

func TestFormatRuleList_Multiple(t *testing.T) {
	result := formatRuleList([]string{"rule.a", "rule.b"})
	if !strings.Contains(result, "rule.a") || !strings.Contains(result, "rule.b") {
		t.Errorf("expected both rules in output, got %q", result)
	}
}

func TestValueOrNone_Empty(t *testing.T) {
	if valueOrNone("") != "<none>" {
		t.Error("expected <none> for empty string")
	}
}

func TestValueOrNone_NonEmpty(t *testing.T) {
	if valueOrNone("hello") != "hello" {
		t.Error("expected original value")
	}
}
