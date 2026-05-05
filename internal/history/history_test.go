package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func samplePort() scanner.Port {
	return scanner.Port{Proto: "tcp", Addr: "0.0.0.0", Port: 8080}
}

func TestRecordAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.jsonl")
	store := history.NewStore(path)

	now := time.Now().UTC().Truncate(time.Second)
	if err := store.Record("added", samplePort(), now); err != nil {
		t.Fatalf("Record: %v", err)
	}

	events, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "added" {
		t.Errorf("expected kind=added, got %s", events[0].Kind)
	}
	if events[0].Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Port.Port)
	}
}

func TestLoadNoFile(t *testing.T) {
	store := history.NewStore("/tmp/portwatch-nonexistent-history.jsonl")
	events, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty slice, got %d events", len(events))
	}
}

func TestRecordMultipleEvents(t *testing.T) {
	dir := t.TempDir()
	store := history.NewStore(filepath.Join(dir, "h.jsonl"))
	now := time.Now().UTC()

	_ = store.Record("added", samplePort(), now)
	_ = store.Record("removed", samplePort(), now.Add(time.Minute))

	events, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[1].Kind != "removed" {
		t.Errorf("expected second event kind=removed")
	}
}

func TestRecordCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "history.jsonl")
	store := history.NewStore(path)

	if err := store.Record("added", samplePort(), time.Now()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
