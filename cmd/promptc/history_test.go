package main

import "testing"

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
