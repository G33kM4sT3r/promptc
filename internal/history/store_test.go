package history

import (
	"fmt"
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

func entry(input string) Entry {
	return Entry{
		Input:     input,
		Timestamp: time.Now(),
	}
}

func TestStoreDelete(t *testing.T) {
	store := NewStore(t.TempDir())
	_ = store.Add(entry("first"))
	_ = store.Add(entry("second"))
	_ = store.Add(entry("third"))

	if err := store.Delete(1); err != nil {
		t.Fatal(err)
	}

	entries, _ := store.List()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Input != "first" || entries[1].Input != "third" {
		t.Error("wrong entries remaining after delete")
	}
}

func TestStoreDeleteOutOfRange(t *testing.T) {
	store := NewStore(t.TempDir())
	_ = store.Add(entry("only"))

	if err := store.Delete(5); err == nil {
		t.Error("expected error for out-of-range delete")
	}
}

func TestStoreClear(t *testing.T) {
	store := NewStore(t.TempDir())
	_ = store.Add(entry("first"))
	_ = store.Add(entry("second"))

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	entries, _ := store.List()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestStoreSearch(t *testing.T) {
	store := NewStore(t.TempDir())
	_ = store.Add(entry("explain closures"))
	_ = store.Add(entry("generate REST API"))
	_ = store.Add(entry("explain generics"))

	results := store.Search("explain")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	results = store.Search("CLOSURES") // case insensitive
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}

	results = store.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestStoreLen(t *testing.T) {
	store := NewStore(t.TempDir())
	n, _ := store.Len()
	if n != 0 {
		t.Errorf("empty store Len = %d, want 0", n)
	}
	_ = store.Add(entry("first"))
	_ = store.Add(entry("second"))
	n, _ = store.Len()
	if n != 2 {
		t.Errorf("Len = %d, want 2", n)
	}
}

func TestStorePrune(t *testing.T) {
	store := NewStore(t.TempDir())
	for i := 0; i < 10; i++ {
		_ = store.Add(entry(fmt.Sprintf("entry-%d", i)))
	}
	entries, _ := store.List()
	if len(entries) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(entries))
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
