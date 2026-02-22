package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"promptc/internal/prompt"
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

// Add appends an entry to the history.
func (s *Store) Add(entry Entry) error {
	entries, _ := s.List() // ignore error for new file
	entries = append(entries, entry)
	return s.save(entries)
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

func (s *Store) save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling history: %w", err)
	}
	return os.WriteFile(s.path, data, 0644)
}
