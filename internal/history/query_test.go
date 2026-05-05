package history_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func seedStore(t *testing.T) *history.Store {
	t.Helper()
	dir := t.TempDir()
	store := history.NewStore(filepath.Join(dir, "h.jsonl"))

	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	_ = store.Record("added", scanner.Port{Proto: "tcp", Addr: "0.0.0.0", Port: 80}, base)
	_ = store.Record("removed", scanner.Port{Proto: "tcp", Addr: "0.0.0.0", Port: 80}, base.Add(time.Hour))
	_ = store.Record("added", scanner.Port{Proto: "udp", Addr: "0.0.0.0", Port: 53}, base.Add(2*time.Hour))
	return store
}

func TestQueryByKind(t *testing.T) {
	store := seedStore(t)
	events, err := store.Query(history.Filter{Kind: "added"})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 added events, got %d", len(events))
	}
}

func TestQueryByProto(t *testing.T) {
	store := seedStore(t)
	events, err := store.Query(history.Filter{Proto: "udp"})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 udp event, got %d", len(events))
	}
}

func TestQuerySince(t *testing.T) {
	store := seedStore(t)
	base := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
	events, err := store.Query(history.Filter{Since: base})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event after cutoff, got %d", len(events))
	}
}

func TestQueryNoFilter(t *testing.T) {
	store := seedStore(t)
	events, err := store.Query(history.Filter{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}
