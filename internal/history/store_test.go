package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"promptc/internal/prompt"
)

func TestStoreAddAndList(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	entry := Entry{
		Input:     "explain closures",
		Spec:      prompt.PromptSpec{Objective: "Explain closures"},
		Score:     72,
		Language:  "en",
		Timestamp: time.Now(),
	}

	if err := store.Add(entry); err != nil {
		t.Fatalf("Add: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].Input != "explain closures" {
		t.Errorf("input = %q, want %q", entries[0].Input, "explain closures")
	}
}

func TestStoreGet(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	_ = store.Add(Entry{Input: "first", Timestamp: time.Now()})
	_ = store.Add(Entry{Input: "second", Timestamp: time.Now()})

	entry, err := store.Get(0)
	if err != nil {
		t.Fatalf("Get(0): %v", err)
	}
	if entry.Input != "first" {
		t.Errorf("entry 0 = %q, want %q", entry.Input, "first")
	}

	entry, err = store.Get(1)
	if err != nil {
		t.Fatalf("Get(1): %v", err)
	}
	if entry.Input != "second" {
		t.Errorf("entry 1 = %q, want %q", entry.Input, "second")
	}
}

func TestStoreGetOutOfRange(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	_, err := store.Get(0)
	if err == nil {
		t.Error("expected error for empty store")
	}
}

func TestStorePersistence(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	_ = store.Add(Entry{Input: "persisted", Timestamp: time.Now()})

	// New store instance reads same file
	store2 := NewStore(dir)
	entries, err := store2.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 || entries[0].Input != "persisted" {
		t.Errorf("persistence failed: got %v", entries)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(dir, "history.json")); os.IsNotExist(err) {
		t.Error("history.json not created")
	}
}
