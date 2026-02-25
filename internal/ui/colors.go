// Package ui provides shared color palette, gradient text rendering,
// NO_COLOR/TTY detection, and a reusable bubbletea spinner for CLI output.
package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var colorsOff bool

// Section header styles — matches existing repl/style.go palette.
var (
	StyleRole        = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	StyleObjective   = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
	StyleContext     = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	StyleScope       = lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Bold(true)
	StyleConstraints = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	StyleOutput      = lipgloss.NewStyle().Foreground(lipgloss.Color("171")).Bold(true)
	StyleQuality     = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)

	// Utility styles for general UI elements.
	StyleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true)
	StyleBold    = lipgloss.NewStyle().Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	StyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)
	StyleContent = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	StyleBullet  = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	StyleTrace   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	StyleSep     = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

	// SectionStyles maps translation label keys to their lipgloss styles.
	SectionStyles = map[string]lipgloss.Style{
		"labels.role":        StyleRole,
		"labels.objective":   StyleObjective,
		"labels.context":     StyleContext,
		"labels.scope":       StyleScope,
		"labels.constraints": StyleConstraints,
		"labels.output":      StyleOutput,
		"labels.quality":     StyleQuality,
	}
)

func init() {
	initColorState()
}

// initColorState detects NO_COLOR env and TTY status.
func initColorState() {
	_, noColor := os.LookupEnv("NO_COLOR")
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	colorsOff = noColor || !isTTY
}

// ColorsDisabled reports whether color output should be suppressed.
func ColorsDisabled() bool {
	return colorsOff
}

// GradientText renders multi-line text with a horizontal color gradient.
// Each line is independently gradient-colored from startHex to endHex.
// Returns plain text if colors are disabled.
func GradientText(text, startHex, endHex string) string {
	if colorsOff {
		return text
	}
	lines := strings.Split(text, "\n")
	var out strings.Builder
	for i, line := range lines {
		if i > 0 {
			out.WriteString("\n")
		}
		runes := []rune(line)
		if len(runes) == 0 {
			continue
		}
		sr, sg, sb := hexToRGB(startHex)
		er, eg, eb := hexToRGB(endHex)
		n := len(runes)
		for j, r := range runes {
			t := 0.0
			if n > 1 {
				t = float64(j) / float64(n-1)
			}
			cr := sr + int(t*float64(er-sr))
			cg := sg + int(t*float64(eg-sg))
			cb := sb + int(t*float64(eb-sb))
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(rgbToHex(cr, cg, cb)))
			out.WriteString(style.Render(string(r)))
		}
	}
	return out.String()
}

// BarChart renders a horizontal bar of the given width.
// Filled portion is proportional to score/max.
func BarChart(score, maxVal, width int) string {
	if maxVal <= 0 {
		maxVal = 1
	}
	filled := 0
	if score > 0 {
		filled = (score * width) / maxVal
		if filled > width {
			filled = width
		}
	}
	empty := width - filled

	fillStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)

	if colorsOff {
		return fillStr + emptyStr
	}

	fillStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	return fillStyle.Render(fillStr) + emptyStyle.Render(emptyStr)
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	r := hexByte(hex[0:2])
	g := hexByte(hex[2:4])
	b := hexByte(hex[4:6])
	return r, g, b
}

func hexByte(s string) int {
	v := 0
	for _, c := range s {
		v *= 16
		if c >= '0' && c <= '9' {
			v += int(c - '0')
		} else if c >= 'a' && c <= 'f' {
			v += int(c-'a') + 10
		} else if c >= 'A' && c <= 'F' {
			v += int(c-'A') + 10
		}
	}
	return v
}

func rgbToHex(r, g, b int) string {
	const hex = "0123456789abcdef"
	return "#" + string([]byte{hex[r/16], hex[r%16], hex[g/16], hex[g%16], hex[b/16], hex[b%16]})
}
