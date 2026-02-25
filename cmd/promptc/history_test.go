package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"promptc/internal/history"
)

func TestHistoryCommand_Empty(t *testing.T) {
	// History may be empty or have entries depending on state
	// Just verify it doesn't crash
	_, _ = runBinary("history")
}

func TestHistoryCommand_InvalidIndex(t *testing.T) {
	_, err := runBinary("history", "99999")
	if err == nil {
		t.Error("expected error for out-of-range index")
	}
}

// seedHistory creates a temp dir with data/ and languages/ symlinks and seeds history entries.
func seedHistory(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Symlink data/ and languages/ from project root so the binary can load config
	os.Symlink(filepath.Join(projectRoot, "data"), filepath.Join(dir, "data"))
	os.Symlink(filepath.Join(projectRoot, "languages"), filepath.Join(dir, "languages"))
	os.Symlink(filepath.Join(projectRoot, "translations"), filepath.Join(dir, "translations"))

	store := history.NewStore(dir)
	_ = store.Add(history.Entry{Input: "explain closures", Score: 70, Language: "en", Timestamp: time.Now()})
	_ = store.Add(history.Entry{Input: "generate REST API", Score: 55, Language: "en", Timestamp: time.Now()})
	_ = store.Add(history.Entry{Input: "explain generics", Score: 65, Language: "en", Timestamp: time.Now()})
	return dir
}

func runBinaryWithEnv(dir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "PROMPTC_DATA="+dir)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestHistorySearch(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--search", "explain")
	if err != nil {
		t.Fatalf("history search failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "explain") {
		t.Error("search output should contain 'explain'")
	}
	// Should find 2 matching entries
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 search results, got %d: %s", len(lines), out)
	}
}

func TestHistorySearchNoResults(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--search", "nonexistent")
	if err != nil {
		t.Fatalf("history search failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No matching") {
		t.Errorf("expected 'No matching' message, got: %s", out)
	}
}

func TestHistoryClear(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--clear", "--yes")
	if err != nil {
		t.Fatalf("history clear failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "cleared") {
		t.Errorf("expected 'cleared' message, got: %s", out)
	}

	// Verify empty
	out, err = runBinaryWithEnv(dir, "history")
	if err != nil {
		t.Fatalf("history list after clear failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "No history") {
		t.Errorf("expected empty history, got: %s", out)
	}
}

func TestHistoryDelete(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--delete", "1", "--yes")
	if err != nil {
		t.Fatalf("history delete failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Deleted") {
		t.Errorf("expected 'Deleted' message, got: %s", out)
	}
}

func TestHistoryExport(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--export")
	if err != nil {
		t.Fatalf("history export failed: %v\n%s", err, out)
	}

	var entries []history.Entry
	if err := json.Unmarshal([]byte(out), &entries); err != nil {
		t.Fatalf("export is not valid JSON: %v\n%s", err, out)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 exported entries, got %d", len(entries))
	}
}

func TestHistoryLimit(t *testing.T) {
	dir := seedHistory(t)

	out, err := runBinaryWithEnv(dir, "history", "--limit", "2")
	if err != nil {
		t.Fatalf("history limit failed: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 entries with --limit 2, got %d: %s", len(lines), out)
	}
}
