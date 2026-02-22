package repl

import (
	"github.com/charmbracelet/lipgloss"
	"promptc/internal/ui"
)

var (
	// REPL-specific styles (not shared)
	stylePrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	styleStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Background(lipgloss.Color("236"))
	styleBanner = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)

	// Aliases to shared styles from ui package
	styleContent = ui.StyleContent
	styleBullet  = ui.StyleBullet
	styleDim     = ui.StyleDim
	styleError   = ui.StyleError
	styleTrace   = ui.StyleTrace
	styleSep     = ui.StyleSep

	// Shared section styles map
	sectionStyles = ui.SectionStyles
)
