package render

import "promptc/internal/prompt"

type Renderer interface {
	Render(prompt.PromptSpec) string
}
