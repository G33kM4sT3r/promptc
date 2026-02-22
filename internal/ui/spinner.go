package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerResult is sent when the work function completes.
type SpinnerResult struct {
	Output string
	Err    error
}

type spinnerModel struct {
	spinner  spinner.Model
	message  string
	result   *SpinnerResult
	quitting bool
}

// RunWithSpinner displays a spinner with the given message while workFn runs.
// Returns the output string and error from workFn.
// If colors are disabled, runs workFn directly without a spinner.
func RunWithSpinner(message string, workFn func() (string, error)) (string, error) {
	if ColorsDisabled() {
		return workFn()
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

	m := spinnerModel{
		spinner: s,
		message: message,
	}

	p := tea.NewProgram(m)

	go func() {
		out, err := workFn()
		p.Send(SpinnerResult{Output: out, Err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	fm, ok := finalModel.(spinnerModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type")
	}
	if fm.result != nil {
		return fm.result.Output, fm.result.Err
	}
	return "", fmt.Errorf("spinner exited without result")
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SpinnerResult:
		m.result = &msg
		m.quitting = true
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return m.spinner.View() + " " + StyleDim.Render(m.message) + "\n"
}
