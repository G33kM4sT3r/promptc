package ui

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestColorsDisabled_WhenNO_COLOR_Set(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	initColorState()
	if !ColorsDisabled() {
		t.Error("expected colors disabled when NO_COLOR is set")
	}
}

func TestColorsDisabled_WhenNO_COLOR_Unset(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	_ = os.Unsetenv("NO_COLOR")
	initColorState()
	_ = ColorsDisabled()
}

func TestGradientText_ReturnsNonEmpty(t *testing.T) {
	oldState := colorsOff
	colorsOff = false
	defer func() { colorsOff = oldState }()

	result := GradientText("HELLO", "#ff00ff", "#00ffff")
	if result == "" {
		t.Error("expected non-empty gradient text")
	}
}

func TestGradientText_PlainWhenDisabled(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	initColorState()
	result := GradientText("HELLO", "#ff00ff", "#00ffff")
	if result != "HELLO" {
		t.Errorf("expected plain text when colors disabled, got %q", result)
	}
}

func TestGradientText_MultiLine(t *testing.T) {
	oldState := colorsOff
	colorsOff = false
	defer func() { colorsOff = oldState }()

	result := GradientText("line1\nline2", "#FF0000", "#0000FF")
	if result == "" {
		t.Error("expected non-empty gradient text")
	}
	// Result contains both lines (may or may not have ANSI escapes depending on lipgloss profile)
	if len(result) < len("line1\nline2") {
		t.Error("expected result at least as long as input")
	}
}

func TestGradientText_SingleChar(t *testing.T) {
	oldState := colorsOff
	colorsOff = false
	defer func() { colorsOff = oldState }()

	result := GradientText("A", "#FF0000", "#0000FF")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestGradientText_EmptyString(t *testing.T) {
	result := GradientText("", "#FF0000", "#0000FF")
	if result != "" {
		t.Errorf("expected empty result for empty input, got %q", result)
	}
}

func TestGradientText_EmptyLineInMiddle(t *testing.T) {
	oldState := colorsOff
	colorsOff = false
	defer func() { colorsOff = oldState }()

	result := GradientText("line1\n\nline3", "#FF0000", "#0000FF")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestBarChart(t *testing.T) {
	oldState := colorsOff
	colorsOff = true
	defer func() { colorsOff = oldState }()

	tests := []struct {
		name     string
		score    int
		max      int
		width    int
		wantFill int
	}{
		{"full", 25, 25, 20, 20},
		{"empty", 0, 15, 20, 0},
		{"half", 5, 10, 20, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BarChart(tt.score, tt.max, tt.width)
			if got == "" {
				t.Error("BarChart returned empty string")
			}
			fills := strings.Count(got, "█")
			empties := strings.Count(got, "░")
			if fills+empties != tt.width {
				t.Errorf("bar width = %d, want %d", fills+empties, tt.width)
			}
			if fills != tt.wantFill {
				t.Errorf("filled = %d, want %d", fills, tt.wantFill)
			}
		})
	}
}

func TestBarChart_ZeroMax(t *testing.T) {
	oldState := colorsOff
	colorsOff = true
	defer func() { colorsOff = oldState }()

	got := BarChart(5, 0, 10)
	if got == "" {
		t.Error("BarChart returned empty string")
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex     string
		r, g, b int
	}{
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
		{"#ffffff", 255, 255, 255},
		{"#000000", 0, 0, 0},
	}
	for _, tt := range tests {
		r, g, b := hexToRGB(tt.hex)
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("hexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)", tt.hex, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestHexToRGB_InvalidLength(t *testing.T) {
	r, g, b := hexToRGB("abc")
	if r != 255 || g != 255 || b != 255 {
		t.Errorf("hexToRGB(invalid) = (%d,%d,%d), want (255,255,255)", r, g, b)
	}
}

func TestRgbToHex(t *testing.T) {
	tests := []struct {
		r, g, b int
		want    string
	}{
		{255, 0, 0, "#ff0000"},
		{0, 255, 0, "#00ff00"},
		{0, 0, 255, "#0000ff"},
	}
	for _, tt := range tests {
		got := rgbToHex(tt.r, tt.g, tt.b)
		if got != tt.want {
			t.Errorf("rgbToHex(%d,%d,%d) = %q, want %q", tt.r, tt.g, tt.b, got, tt.want)
		}
	}
}

func TestHexByte(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"00", 0},
		{"ff", 255},
		{"FF", 255},
		{"0a", 10},
		{"7f", 127},
	}
	for _, tt := range tests {
		got := hexByte(tt.input)
		if got != tt.want {
			t.Errorf("hexByte(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestSpinnerResult_Type(t *testing.T) {
	r := SpinnerResult{Output: "test", Err: nil}
	if r.Output != "test" {
		t.Error("expected output to be 'test'")
	}
	if r.Err != nil {
		t.Error("expected nil error")
	}

	r2 := SpinnerResult{Output: "", Err: errors.New("fail")}
	if r2.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRunWithSpinner_ColorsDisabled(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	initColorState()

	out, err := RunWithSpinner("loading...", func() (string, error) {
		return "result", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "result" {
		t.Errorf("expected 'result', got %q", out)
	}
}

func TestRunWithSpinner_ColorsDisabled_Error(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	initColorState()

	_, err := RunWithSpinner("loading...", func() (string, error) {
		return "", errors.New("work failed")
	})
	if err == nil || err.Error() != "work failed" {
		t.Errorf("expected 'work failed' error, got %v", err)
	}
}

func TestSpinnerModel_Init(t *testing.T) {
	m := spinnerModel{
		spinner: spinner.New(),
		message: "loading",
	}
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected non-nil init command")
	}
}

func TestSpinnerModel_Update_SpinnerResult(t *testing.T) {
	m := spinnerModel{
		spinner: spinner.New(),
		message: "loading",
	}
	result := SpinnerResult{Output: "done", Err: nil}
	newModel, cmd := m.Update(result)
	sm, ok := newModel.(spinnerModel)
	if !ok {
		t.Fatal("expected spinnerModel type")
	}
	if !sm.quitting {
		t.Error("expected quitting to be true after SpinnerResult")
	}
	if sm.result == nil || sm.result.Output != "done" {
		t.Error("expected result to be stored")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command")
	}
}

func TestSpinnerModel_Update_CtrlC(t *testing.T) {
	m := spinnerModel{
		spinner: spinner.New(),
		message: "loading",
	}
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := m.Update(msg)
	sm, ok := newModel.(spinnerModel)
	if !ok {
		t.Fatal("expected spinnerModel type")
	}
	if !sm.quitting {
		t.Error("expected quitting after Ctrl+C")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command")
	}
}

func TestSpinnerModel_Update_RegularKey(t *testing.T) {
	m := spinnerModel{
		spinner: spinner.New(),
		message: "loading",
	}
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ := m.Update(msg)
	sm, ok := newModel.(spinnerModel)
	if !ok {
		t.Fatal("expected spinnerModel type")
	}
	if sm.quitting {
		t.Error("regular key should not cause quitting")
	}
}

func TestSpinnerModel_View(t *testing.T) {
	m := spinnerModel{
		spinner: spinner.New(),
		message: "loading",
	}
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}

	m.quitting = true
	view = m.View()
	if view != "" {
		t.Errorf("expected empty view when quitting, got %q", view)
	}
}
