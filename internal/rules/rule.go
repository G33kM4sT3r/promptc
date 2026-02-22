package rules

import (
	"promptc/internal/i18n"
	"promptc/internal/prompt"
	"promptc/internal/slots"
)

type Rule struct {
	ID    string
	When  func(slots.Slots) bool
	Apply func(*prompt.PromptSpec, slots.Slots, *i18n.Translator)
}
