package repl

import "testing"

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantTyp commandType
		wantArg string
	}{
		// Valid quit commands
		{"quit :q", ":q", cmdQuit, ""},
		{"quit :quit", ":quit", cmdQuit, ""},
		{"quit :exit", ":exit", cmdQuit, ""},

		// Explain
		{"explain", ":explain", cmdExplain, ""},

		// Help
		{"help :help", ":help", cmdHelp, ""},
		{"help :h", ":h", cmdHelp, ""},

		// Lang with argument
		{"lang de", ":lang de", cmdLang, "de"},
		{"language en", ":language en", cmdLang, "en"},

		// Lang with no argument
		{"lang no arg", ":lang", cmdLang, ""},
		{"language no arg", ":language", cmdLang, ""},

		// Case insensitivity
		{"case :Q", ":Q", cmdQuit, ""},
		{"case :QUIT", ":QUIT", cmdQuit, ""},
		{"case :Help", ":Help", cmdHelp, ""},
		{"case :LANG de", ":LANG de", cmdLang, "de"},
		{"case :Explain", ":Explain", cmdExplain, ""},

		// Whitespace handling
		{"whitespace quit", "  :quit  ", cmdQuit, ""},
		{"whitespace lang", ":lang  fr ", cmdLang, "fr"},
		{"whitespace leading", "  :help", cmdHelp, ""},

		// Empty and regular input (not commands)
		{"empty", "", cmdNone, ""},
		{"regular text", "explain docker", cmdNone, ""},
		{"no colon prefix", "quit", cmdNone, ""},

		// Unknown commands
		{"unknown :foo", ":foo", cmdUnknown, ":foo"},
		{"unknown :bar baz", ":bar baz", cmdUnknown, ":bar baz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.input)
			if got.typ != tt.wantTyp {
				t.Errorf("parseCommand(%q).typ = %d, want %d", tt.input, got.typ, tt.wantTyp)
			}
			if got.args != tt.wantArg {
				t.Errorf("parseCommand(%q).args = %q, want %q", tt.input, got.args, tt.wantArg)
			}
		})
	}
}
