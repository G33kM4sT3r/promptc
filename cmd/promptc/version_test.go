package main

import (
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	out, err := runBinary("version")
	if err != nil {
		t.Fatalf("version command failed: %v\n%s", err, out)
	}
	if !strings.Contains(strings.ToLower(out), "promptc") {
		t.Errorf("expected version info, got: %s", out)
	}
}
