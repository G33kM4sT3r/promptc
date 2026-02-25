package repl

import (
	"strings"
)

type commandType int

const (
	cmdNone commandType = iota
	cmdQuit
	cmdExplain
	cmdLang
	cmdHelp
	cmdOutput
	cmdHistory
	cmdRecall
	cmdCopy
	cmdSearch
	cmdUnknown
)

type command struct {
	typ  commandType
	args string
}

// parseCommand parses REPL commands from user input.
// Commands start with ':' and are case insensitive.
func parseCommand(input string) command {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, ":") {
		return command{typ: cmdNone}
	}

	lower := strings.ToLower(input)
	parts := strings.SplitN(lower, " ", 2)
	cmd := parts[0]
	args := ""
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}

	switch cmd {
	case ":q", ":quit", ":exit":
		return command{typ: cmdQuit}
	case ":explain":
		return command{typ: cmdExplain}
	case ":lang", ":language":
		return command{typ: cmdLang, args: args}
	case ":help", ":h":
		return command{typ: cmdHelp}
	case ":output":
		return command{typ: cmdOutput, args: args}
	case ":history":
		return command{typ: cmdHistory}
	case ":recall":
		return command{typ: cmdRecall, args: args}
	case ":copy":
		return command{typ: cmdCopy}
	case ":search":
		return command{typ: cmdSearch, args: args}
	default:
		return command{typ: cmdUnknown, args: input}
	}
}

// helpText returns styled help output for the REPL.
func helpText() string {
	header := styleBanner.Render("Commands:")
	lines := []string{
		header,
		"  " + stylePrompt.Render(":help, :h") + styleDim.Render("          Show this help"),
		"  " + stylePrompt.Render(":explain") + styleDim.Render("           Toggle explain mode"),
		"  " + stylePrompt.Render(":lang <code>") + styleDim.Render("       Set language (e.g., :lang de)"),
		"  " + stylePrompt.Render(":lang") + styleDim.Render("              Show current language"),
		"  " + stylePrompt.Render(":output <fmt>") + styleDim.Render("       Set output format (text, json, yaml)"),
		"  " + stylePrompt.Render(":history") + styleDim.Render("           Show recent history"),
		"  " + stylePrompt.Render(":recall <N>") + styleDim.Render("        Recall history entry by index"),
		"  " + stylePrompt.Render(":search <term>") + styleDim.Render("     Search history by input text"),
		"  " + stylePrompt.Render(":copy") + styleDim.Render("              Copy last output to clipboard"),
		"  " + stylePrompt.Render(":quit, :q, :exit") + styleDim.Render("   Exit the REPL"),
		"",
		styleDim.Render("  Type any text to compile it into a structured prompt."),
	}
	return strings.Join(lines, "\n")
}
