package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"promptc/internal/prompt"
)

const (
	defaultMaxEntries = 500
	defaultMaxAge     = 90 * 24 * time.Hour
)

// Entry represents a single prompt compilation in the history.
type Entry struct {
	Input     string            `json:"input"`
	Spec      prompt.PromptSpec `json:"spec"`
	Score     int               `json:"score"`
	Language  string            `json:"language"`
	Timestamp time.Time         `json:"timestamp"`
}

// Store manages prompt compilation history on disk.
type Store struct {
	path string
}

// NewStore creates a history store that persists to history.json in the given directory.
func NewStore(dir string) *Store {
	return &Store{path: filepath.Join(dir, "history.json")}
}

// Add appends an entry to the history and auto-prunes old/excess entries.
func (s *Store) Add(entry Entry) error {
	entries, _ := s.List() // ignore error for new file
	entries = append(entries, entry)
	if err := s.save(entries); err != nil {
		return err
	}
	return s.Prune()
}

// List returns all history entries.
func (s *Store) List() ([]Entry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading history: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing history: %w", err)
	}
	return entries, nil
}

// Get returns a specific history entry by index.
func (s *Store) Get(index int) (Entry, error) {
	entries, err := s.List()
	if err != nil {
		return Entry{}, err
	}
	if index < 0 || index >= len(entries) {
		return Entry{}, fmt.Errorf("index %d out of range (history has %d entries)", index, len(entries))
	}
	return entries[index], nil
}

// Delete removes the entry at the given index.
func (s *Store) Delete(index int) error {
	entries, err := s.List()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(entries) {
		return fmt.Errorf("index %d out of range (history has %d entries)", index, len(entries))
	}
	entries = append(entries[:index], entries[index+1:]...)
	return s.save(entries)
}

// Clear removes all history entries.
func (s *Store) Clear() error {
	return s.save([]Entry{})
}

// Search returns entries whose input contains the term (case-insensitive).
func (s *Store) Search(term string) []Entry {
	entries, err := s.List()
	if err != nil {
		return nil
	}
	lower := strings.ToLower(term)
	var results []Entry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Input), lower) {
			results = append(results, e)
		}
	}
	return results
}

// Len returns the number of history entries.
func (s *Store) Len() (int, error) {
	entries, err := s.List()
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// Prune removes entries exceeding the max count and max age limits.
func (s *Store) Prune() error {
	entries, err := s.List()
	if err != nil {
		return err
	}

	now := time.Now()
	var pruned []Entry
	for _, e := range entries {
		if now.Sub(e.Timestamp) <= defaultMaxAge {
			pruned = append(pruned, e)
		}
	}

	if len(pruned) > defaultMaxEntries {
		pruned = pruned[len(pruned)-defaultMaxEntries:]
	}

	return s.save(pruned)
}

func (s *Store) save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling history: %w", err)
	}
	return os.WriteFile(s.path, data, 0644)
}
