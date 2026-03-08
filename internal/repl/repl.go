package repl

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"promptc/internal/clipboard"
	"promptc/internal/config"
	"promptc/internal/history"
	"promptc/internal/i18n"
	"promptc/internal/language"
	"promptc/internal/pipeline"
	"promptc/internal/prompt"
	"promptc/internal/render"
	"promptc/internal/score"
)

type model struct {
	textInput  textinput.Model
	cfg        *config.Config
	detector   *language.Detector
	baseDir    string
	lang       string
	explainOn  bool
	outputFmt  string
	store      *history.Store
	lastOutput string
	output     string
	err        error
	quitting   bool
	version    string
}

// Run starts the interactive REPL.
func Run(cfg *config.Config, detector *language.Detector, baseDir, lang, version string) error {
	ti := textinput.New()
	ti.Prompt = stylePrompt.Render("promptc> ")
	ti.Focus()

	m := model{
		textInput: ti,
		cfg:       cfg,
		detector:  detector,
		baseDir:   baseDir,
		lang:      lang,
		outputFmt: "text",
		store:     history.NewStore(baseDir),
		version:   version,
	}

	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(m.textInput.Value())
			m.textInput.SetValue("")

			if input == "" {
				return m, nil
			}

			cmd := parseCommand(input)
			switch cmd.typ {
			case cmdQuit:
				m.quitting = true
				return m, tea.Quit
			case cmdExplain:
				m.explainOn = !m.explainOn
				state := "off"
				if m.explainOn {
					state = "on"
				}
				m.output = styleDim.Render(fmt.Sprintf("Explain mode: %s", state))
				m.err = nil
				return m, nil
			case cmdLang:
				if cmd.args == "" {
					current := m.lang
					if current == "" {
						current = "auto"
					}
					m.output = styleDim.Render(fmt.Sprintf("Language: %s", current))
				} else {
					m.lang = cmd.args
					m.output = styleDim.Render(fmt.Sprintf("Language set to: %s", cmd.args))
				}
				m.err = nil
				return m, nil
			case cmdOutput:
				switch cmd.args {
				case "text", "json", "yaml":
					m.outputFmt = cmd.args
					m.output = styleDim.Render(fmt.Sprintf("Output format: %s", cmd.args))
				case "":
					m.output = styleDim.Render(fmt.Sprintf("Output format: %s", m.outputFmt))
				default:
					m.output = styleError.Render(fmt.Sprintf("Unknown format: %s (use text, json, yaml)", cmd.args))
				}
				m.err = nil
				return m, nil
			case cmdHistory:
				entries, err := m.store.List()
				if err != nil {
					m.output = styleError.Render(fmt.Sprintf("History error: %v", err))
				} else if len(entries) == 0 {
					m.output = styleDim.Render("No history entries.")
				} else {
					tr := loadTranslator(m.baseDir, m.lang)
					var lines []string
					start := 0
					if len(entries) > 10 {
						start = len(entries) - 10
					}
					for i := start; i < len(entries); i++ {
						e := entries[i]
						ts := tr.FormatTime(e.Timestamp)
						lines = append(lines, styleDim.Render(fmt.Sprintf("[%d] %s  %s  (score: %d)", i, ts, e.Input, e.Score)))
					}
					m.output = strings.Join(lines, "\n")
				}
				m.err = nil
				return m, nil
			case cmdRecall:
				if cmd.args == "" {
					m.output = styleError.Render("Usage: :recall <index>")
					m.err = nil
					return m, nil
				}
				index, err := strconv.Atoi(cmd.args)
				if err != nil {
					m.output = styleError.Render(fmt.Sprintf("Invalid index: %s", cmd.args))
					m.err = nil
					return m, nil
				}
				entry, err := m.store.Get(index)
				if err != nil {
					m.output = styleError.Render(err.Error())
					m.err = nil
					return m, nil
				}
				translator := loadTranslator(m.baseDir, entry.Language)
				rendered := m.renderSpec(entry.Spec, translator)
				scoreLine := styleDim.Render(fmt.Sprintf("Score: %d/100", entry.Score))
				m.output = rendered + "\n" + scoreLine
				m.err = nil
				return m, nil
			case cmdSearch:
				if cmd.args == "" {
					m.output = styleError.Render("Usage: :search <term>")
					m.err = nil
					return m, nil
				}
				results := m.store.Search(cmd.args)
				if len(results) == 0 {
					m.output = styleDim.Render("No matching entries.")
				} else {
					tr := loadTranslator(m.baseDir, m.lang)
					var lines []string
					for i, e := range results {
						ts := tr.FormatTime(e.Timestamp)
						lines = append(lines, styleDim.Render(fmt.Sprintf("[%d] %s  %s  (score: %d)", i, ts, e.Input, e.Score)))
					}
					m.output = strings.Join(lines, "\n")
				}
				m.err = nil
				return m, nil
			case cmdCopy:
				if m.lastOutput == "" {
					m.output = styleError.Render("Nothing to copy.")
				} else if err := clipboard.Copy(m.lastOutput); err != nil {
					m.output = styleError.Render(fmt.Sprintf("Clipboard error: %v", err))
				} else {
					m.output = styleDim.Render("Copied to clipboard.")
				}
				m.err = nil
				return m, nil
			case cmdHelp:
				m.output = helpText()
				m.err = nil
				return m, nil
			case cmdUnknown:
				m.output = styleError.Render(fmt.Sprintf("Unknown command: %s", cmd.args)) + "\n" + styleDim.Render("Type :help for available commands.")
				m.err = nil
				return m, nil
			}

			// Regular input: process through pipeline
			m.output, m.lastOutput, m.err = m.processInput(input)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return styleDim.Render("Goodbye!") + "\n"
	}

	var b strings.Builder

	if m.output == "" && m.err == nil {
		b.WriteString(styleBanner.Render(fmt.Sprintf("promptc %s", m.version)) + " " + styleDim.Render("interactive mode") + "\n")
		b.WriteString(styleDim.Render("Type :help for commands, Ctrl+C to quit.") + "\n\n")
	}

	if m.err != nil {
		b.WriteString(styleError.Render("Error: "+m.err.Error()) + "\n\n")
	} else if m.output != "" {
		b.WriteString(m.output + "\n\n")
	}

	// Status bar
	langDisplay := m.lang
	if langDisplay == "" {
		langDisplay = "auto"
	}
	explainDisplay := "off"
	if m.explainOn {
		explainDisplay = "on"
	}
	statusText := fmt.Sprintf(" lang: %s \u2502 explain: %s \u2502 output: %s ", langDisplay, explainDisplay, m.outputFmt)
	b.WriteString(styleStatus.Render(statusText) + "\n")

	b.WriteString(m.textInput.View())

	return b.String()
}

// processInput returns (styled output, raw output for clipboard, error).
func (m model) processInput(input string) (string, string, error) {
	s, err := pipeline.Run(input, m.cfg, m.detector)
	if err != nil {
		return "", "", err
	}

	lang := s.Language
	if m.lang != "" {
		lang = m.lang
	}
	translator := loadTranslator(m.baseDir, lang)

	var spec prompt.PromptSpec
	var rendered string

	if m.explainOn {
		result := pipeline.ApplyRulesWithTrace(s, translator, m.cfg.Enrichments)
		spec = result.Spec
		rendered = m.renderSpec(result.Spec, translator)
		if m.outputFmt == "text" {
			trace := renderStyledExplain(result.Applied, result.Skipped)
			rendered = rendered + "\n" + trace
		}
	} else {
		spec = pipeline.ApplyRules(s, translator, m.cfg.Enrichments)
		rendered = m.renderSpec(spec, translator)
	}

	scoreResult := score.Score(spec)

	// Auto-save to history
	_ = m.store.Add(history.Entry{
		Input:     input,
		Spec:      spec,
		Score:     scoreResult.Total,
		Language:  lang,
		Timestamp: time.Now(),
	})

	// Raw output for clipboard (use plain text renderer for text mode)
	var rawOutput string
	switch m.outputFmt {
	case "json":
		rawOutput = (&render.JSONRenderer{}).Render(spec)
	case "yaml":
		rawOutput = (&render.YAMLRenderer{}).Render(spec)
	default:
		rawOutput = render.NewTranslated(translator).Render(spec)
	}

	scoreLine := styleDim.Render(fmt.Sprintf("Score: %d/100", scoreResult.Total))
	return rendered + "\n" + scoreLine, rawOutput, nil
}

func (m model) renderSpec(spec prompt.PromptSpec, translator *i18n.Translator) string {
	switch m.outputFmt {
	case "json":
		return (&render.JSONRenderer{}).Render(spec)
	case "yaml":
		return (&render.YAMLRenderer{}).Render(spec)
	default:
		return renderStyledSpec(spec, translator)
	}
}

func renderStyledSpec(p prompt.PromptSpec, t *i18n.Translator) string {
	var b strings.Builder

	// Role (no label, just the value)
	if p.Role != "" {
		style, ok := sectionStyles["labels.role"]
		if !ok {
			style = lipgloss.NewStyle()
		}
		b.WriteString(style.Render(p.Role) + "\n")
	}

	writeStyledSection(&b, t, "labels.objective", p.Objective)
	writeStyledSection(&b, t, "labels.context", p.Context)
	writeStyledList(&b, t, "labels.scope", p.Scope)
	writeStyledList(&b, t, "labels.constraints", p.Constraints)
	writeStyledList(&b, t, "labels.output", p.OutputSpec)
	writeStyledList(&b, t, "labels.quality", p.QualityCriteria)

	return strings.TrimRight(b.String(), "\n")
}

func writeStyledSection(b *strings.Builder, t *i18n.Translator, labelKey, value string) {
	if value == "" {
		return
	}
	style, ok := sectionStyles[labelKey]
	if !ok {
		style = lipgloss.NewStyle()
	}
	label := t.Get(labelKey)
	b.WriteString("\n" + style.Render(label+":") + "\n")
	b.WriteString(styleContent.Render(value) + "\n")
}

func writeStyledList(b *strings.Builder, t *i18n.Translator, labelKey string, items []string) {
	if len(items) == 0 {
		return
	}
	style, ok := sectionStyles[labelKey]
	if !ok {
		style = lipgloss.NewStyle()
	}
	label := t.Get(labelKey)
	b.WriteString("\n" + style.Render(label+":") + "\n")
	for _, item := range items {
		bullet := styleBullet.Render("  - ")
		b.WriteString(bullet + styleContent.Render(item) + "\n")
	}
}

func renderStyledExplain(applied, skipped []string) string {
	var b strings.Builder
	b.WriteString(styleSep.Render(strings.Repeat("─", 40)) + "\n")
	b.WriteString(styleTrace.Render("Rules applied: "+strings.Join(applied, ", ")) + "\n")
	b.WriteString(styleTrace.Render("Rules skipped: "+strings.Join(skipped, ", ")) + "\n")
	return b.String()
}

func loadTranslator(baseDir, lang string) *i18n.Translator {
	if lang == "" || lang == "unknown" {
		lang = "en"
	}
	langDir := filepath.Join(baseDir, "languages")
	translator, err := i18n.Load(langDir, lang, "en")
	if err != nil {
		translator, _ = i18n.Load(langDir, "en", "en")
	}
	return translator
}
