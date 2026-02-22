package repl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"promptc/internal/config"
	"promptc/internal/history"
	"promptc/internal/language"
	"promptc/internal/prompt"
)

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

func newTestModel(t *testing.T) model {
	t.Helper()
	baseDir := findProjectRoot(t)
	cfg, err := config.Load(baseDir)
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	detector, err2 := language.NewDetector(filepath.Join(baseDir, "data", "lid.176.ftz"), cfg)
	if err2 != nil {
		t.Fatalf("NewDetector: %v", err2)
	}

	ti := textinput.New()
	ti.Prompt = "> "
	ti.Focus()

	return model{
		textInput: ti,
		cfg:       cfg,
		detector:  detector,
		baseDir:   baseDir,
		outputFmt: "text",
		store:     history.NewStore(t.TempDir()),
		version:   "test",
	}
}

// --- Update tests ---

func TestUpdate_CtrlC(t *testing.T) {
	m := newTestModel(t)
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := m.Update(msg)
	nm := newModel.(model)
	if !nm.quitting {
		t.Error("expected quitting after Ctrl+C")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command")
	}
}

func TestUpdate_CtrlD(t *testing.T) {
	m := newTestModel(t)
	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !nm.quitting {
		t.Error("expected quitting after Ctrl+D")
	}
}

func TestUpdate_EmptyEnter(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue("")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	nm := newModel.(model)
	if nm.quitting {
		t.Error("empty enter should not quit")
	}
	if cmd != nil {
		t.Error("empty enter should return nil cmd")
	}
}

func TestUpdate_QuitCommand(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":quit")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(msg)
	nm := newModel.(model)
	if !nm.quitting {
		t.Error("expected quitting after :quit")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command")
	}
}

func TestUpdate_ExplainToggle(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":explain")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !nm.explainOn {
		t.Error("expected explain on after first toggle")
	}
	if !strings.Contains(nm.output, "on") {
		t.Error("expected output to contain 'on'")
	}

	// Toggle again
	nm.textInput.SetValue(":explain")
	newModel2, _ := nm.Update(msg)
	nm2 := newModel2.(model)
	if nm2.explainOn {
		t.Error("expected explain off after second toggle")
	}
}

func TestUpdate_LangNoArgs(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":lang")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "auto") {
		t.Error("expected 'auto' in lang output when no lang set")
	}
}

func TestUpdate_LangSetDe(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":lang de")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if nm.lang != "de" {
		t.Errorf("expected lang 'de', got %q", nm.lang)
	}
}

func TestUpdate_OutputFormat(t *testing.T) {
	m := newTestModel(t)

	// Set to json
	m.textInput.SetValue(":output json")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if nm.outputFmt != "json" {
		t.Errorf("expected json, got %q", nm.outputFmt)
	}

	// Query current
	nm.textInput.SetValue(":output")
	newModel2, _ := nm.Update(msg)
	nm2 := newModel2.(model)
	if !strings.Contains(nm2.output, "json") {
		t.Error("expected output to show 'json'")
	}

	// Unknown format
	nm2.textInput.SetValue(":output xml")
	newModel3, _ := nm2.Update(msg)
	nm3 := newModel3.(model)
	if !strings.Contains(nm3.output, "Unknown format") {
		t.Error("expected 'Unknown format' error")
	}
}

func TestUpdate_Help(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":help")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "Commands") {
		t.Error("expected help to contain 'Commands'")
	}
}

func TestUpdate_UnknownCommand(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":foobar")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "Unknown command") {
		t.Error("expected 'Unknown command' in output")
	}
}

func TestUpdate_History_Empty(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":history")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "No history") {
		t.Error("expected 'No history' for empty store")
	}
}

func TestUpdate_Recall_NoArgs(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":recall")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "Usage") {
		t.Error("expected usage hint")
	}
}

func TestUpdate_Recall_InvalidIndex(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":recall abc")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "Invalid index") {
		t.Error("expected 'Invalid index' error")
	}
}

func TestUpdate_Recall_OutOfRange(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":recall 999")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "out of range") {
		t.Error("expected 'out of range' error")
	}
}

func TestUpdate_Copy_NothingToCopy(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue(":copy")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if !strings.Contains(nm.output, "Nothing to copy") {
		t.Error("expected 'Nothing to copy' error")
	}
}

func TestUpdate_ProcessInput(t *testing.T) {
	m := newTestModel(t)
	m.textInput.SetValue("explain kubernetes")
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	nm := newModel.(model)
	if nm.err != nil {
		t.Fatalf("processInput error: %v", nm.err)
	}
	if nm.output == "" {
		t.Error("expected non-empty output for valid input")
	}
	if nm.lastOutput == "" {
		t.Error("expected non-empty lastOutput for clipboard")
	}
}

func TestUpdate_TextInputForwarding(t *testing.T) {
	m := newTestModel(t)
	// Non-key message should be forwarded to textinput
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	newModel, _ := m.Update(msg)
	_ = newModel.(model) // should not panic
}

// --- View tests ---

func TestView_Initial(t *testing.T) {
	m := newTestModel(t)
	view := m.View()
	if !strings.Contains(view, "promptc") {
		t.Error("expected 'promptc' in initial view")
	}
	if !strings.Contains(view, "lang: auto") {
		t.Error("expected 'lang: auto' in status bar")
	}
}

func TestView_Quitting(t *testing.T) {
	m := newTestModel(t)
	m.quitting = true
	view := m.View()
	if !strings.Contains(view, "Goodbye") {
		t.Error("expected 'Goodbye' when quitting")
	}
}

func TestView_WithOutput(t *testing.T) {
	m := newTestModel(t)
	m.output = "test output"
	view := m.View()
	if !strings.Contains(view, "test output") {
		t.Error("expected output in view")
	}
}

func TestView_WithError(t *testing.T) {
	m := newTestModel(t)
	m.err = fmt.Errorf("test error")
	view := m.View()
	if !strings.Contains(view, "test error") {
		t.Error("expected error in view")
	}
}

func TestView_WithLangAndExplain(t *testing.T) {
	m := newTestModel(t)
	m.lang = "de"
	m.explainOn = true
	view := m.View()
	if !strings.Contains(view, "lang: de") {
		t.Error("expected 'lang: de' in status bar")
	}
	if !strings.Contains(view, "explain: on") {
		t.Error("expected 'explain: on' in status bar")
	}
}

// --- Rendering tests ---

func TestRenderStyledSpec(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "en")

	spec := prompt.PromptSpec{
		Role:            "You are a Go expert.",
		Objective:       "Explain closures in Go.",
		Context:         "Beginner audience.",
		Scope:           []string{"definition", "examples"},
		Constraints:     []string{"use simple language"},
		OutputSpec:      []string{"step-by-step"},
		QualityCriteria: []string{"clear", "concise"},
	}

	result := renderStyledSpec(spec, translator)
	if result == "" {
		t.Error("expected non-empty styled spec")
	}
	if !strings.Contains(result, "Go expert") {
		t.Error("expected role in output")
	}
}

func TestRenderStyledSpec_Empty(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "en")

	spec := prompt.PromptSpec{}
	result := renderStyledSpec(spec, translator)
	if result != "" {
		t.Errorf("expected empty result for empty spec, got %q", result)
	}
}

func TestRenderSpec_JSON(t *testing.T) {
	m := newTestModel(t)
	m.outputFmt = "json"
	translator := loadTranslator(m.baseDir, "en")
	spec := prompt.PromptSpec{Objective: "test"}
	result := m.renderSpec(spec, translator)
	if !strings.HasPrefix(strings.TrimSpace(result), "{") {
		t.Error("expected JSON output")
	}
}

func TestRenderSpec_YAML(t *testing.T) {
	m := newTestModel(t)
	m.outputFmt = "yaml"
	translator := loadTranslator(m.baseDir, "en")
	spec := prompt.PromptSpec{Objective: "test"}
	result := m.renderSpec(spec, translator)
	if !strings.Contains(result, "objective:") {
		t.Error("expected YAML output")
	}
}

func TestRenderStyledExplain(t *testing.T) {
	result := renderStyledExplain(
		[]string{"rule.a", "rule.b"},
		[]string{"rule.c"},
	)
	if !strings.Contains(result, "rule.a") {
		t.Error("expected applied rules in output")
	}
	if !strings.Contains(result, "rule.c") {
		t.Error("expected skipped rules in output")
	}
}

func TestHelpText(t *testing.T) {
	h := helpText()
	if !strings.Contains(h, "Commands") {
		t.Error("expected 'Commands' header")
	}
	if !strings.Contains(h, ":quit") {
		t.Error("expected :quit in help")
	}
}

// --- loadTranslator tests ---

func TestLoadTranslator_English(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "en")
	if translator == nil {
		t.Fatal("expected non-nil translator")
	}
	label := translator.Get("labels.objective")
	if label == "" || label == "labels.objective" {
		t.Error("expected translated label for objective")
	}
}

func TestLoadTranslator_German(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "de")
	if translator == nil {
		t.Fatal("expected non-nil translator")
	}
}

func TestLoadTranslator_EmptyFallsBackToEnglish(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "")
	if translator == nil {
		t.Fatal("expected non-nil translator")
	}
}

func TestLoadTranslator_UnknownFallsBackToEnglish(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "unknown")
	if translator == nil {
		t.Fatal("expected non-nil translator")
	}
}

func TestLoadTranslator_InvalidLangFallsBack(t *testing.T) {
	baseDir := findProjectRoot(t)
	translator := loadTranslator(baseDir, "xx")
	if translator == nil {
		t.Fatal("expected non-nil translator (should fall back to en)")
	}
}

// --- Init test ---

func TestInit(t *testing.T) {
	m := newTestModel(t)
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected non-nil init command (textinput.Blink)")
	}
}
