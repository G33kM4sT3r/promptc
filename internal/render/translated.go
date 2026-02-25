package render

import (
	"fmt"
	"strings"

	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/ui"
)

// TranslatedRenderer renders a PromptSpec using translated labels from the i18n Translator.
type TranslatedRenderer struct {
	t *i18n.Translator
}

// NewTranslated creates a renderer that uses the given Translator for section labels.
func NewTranslated(t *i18n.Translator) *TranslatedRenderer {
	return &TranslatedRenderer{t: t}
}

func (r *TranslatedRenderer) Render(p prompt.PromptSpec) string {
	var b strings.Builder

	writeSection(&b, p.Role)
	writeLabeled(&b, r.t.Get("labels.objective"), p.Objective)
	writeLabeled(&b, r.t.Get("labels.context"), p.Context)
	writeList(&b, r.t.Get("labels.scope"), p.Scope)
	writeList(&b, r.t.Get("labels.constraints"), p.Constraints)
	writeList(&b, r.t.Get("labels.output"), p.OutputSpec)
	writeList(&b, r.t.Get("labels.quality"), p.QualityCriteria)

	return strings.TrimSpace(b.String())
}

// RenderScore renders a score breakdown with bar charts.
func (r *TranslatedRenderer) RenderScore(sb ScoreBreakdown) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\nScore: %d/100\n\n", sb.Total)

	sections := []string{"objective", "context", "scope", "output", "role", "constraints", "quality"}
	for _, sec := range sections {
		maxW := sb.MaxWeights[sec]
		sc := sb.Breakdown[sec]
		bar := ui.BarChart(sc, maxW, 20)
		fmt.Fprintf(&b, "  %-14s %2d/%-3d %s\n", sec, sc, maxW, bar)
	}

	return b.String()
}
