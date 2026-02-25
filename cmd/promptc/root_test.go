package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	binaryPath  string
	projectRoot string
)

func TestMain(m *testing.M) {
	// Find project root (walk up from cwd until go.mod found)
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root")
		}
		dir = parent
	}

	// Build the binary once for all tests
	tmpDir, err := os.MkdirTemp("", "promptc-test")
	if err != nil {
		panic(err)
	}
	binaryPath = filepath.Join(tmpDir, "promptc")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(string(out) + ": " + err.Error())
	}
	code := m.Run()
	_ = os.RemoveAll(tmpDir)
	os.Exit(code)
}

func runBinary(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "PROMPTC_DATA="+projectRoot)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runBinaryWithStdin(stdin string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "PROMPTC_DATA="+projectRoot)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestRootCommand_ShowsWelcome(t *testing.T) {
	out, err := runBinary()
	if err != nil {
		t.Fatalf("root command failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "promptc") {
		t.Errorf("expected welcome to mention promptc, got: %s", out)
	}
}

func TestFormatBuildDate_Valid(t *testing.T) {
	got := formatBuildDate("2026-02-22T10:00:00Z")
	if !strings.Contains(got, "2026") {
		t.Errorf("formatBuildDate valid = %q, expected date containing 2026", got)
	}
}

func TestFormatBuildDate_Invalid(t *testing.T) {
	got := formatBuildDate("invalid")
	if got != "invalid" {
		t.Errorf("formatBuildDate invalid = %q, expected 'invalid'", got)
	}
}

func TestFormatBuildDate_Empty(t *testing.T) {
	got := formatBuildDate("")
	if got != "" {
		t.Errorf("formatBuildDate empty = %q, expected empty", got)
	}
}
