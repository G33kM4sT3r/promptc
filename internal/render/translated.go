package render

import (
	"strings"

	"promptc/internal/i18n"
	"promptc/internal/prompt"
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
