package render

import (
	"strings"
)

// ScoreBreakdown holds data for rendering a score breakdown.
type ScoreBreakdown struct {
	Total      int            `json:"total" yaml:"total"`
	Breakdown  map[string]int `json:"breakdown" yaml:"breakdown"`
	MaxWeights map[string]int `json:"max_weights" yaml:"max_weights"`
}

func writeSection(b *strings.Builder, value string) {
	if value == "" {
		return
	}
	b.WriteString(value)
	b.WriteString("\n\n")
}

func writeLabeled(b *strings.Builder, label, value string) {
	if value == "" {
		return
	}
	b.WriteString(label)
	b.WriteString(":\n")
	b.WriteString(value)
	b.WriteString("\n\n")
}

func writeList(b *strings.Builder, label string, items []string) {
	if len(items) == 0 {
		return
	}
	b.WriteString(label)
	b.WriteString(":\n")
	for _, item := range items {
		b.WriteString("- ")
		b.WriteString(strings.TrimSpace(item))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}
