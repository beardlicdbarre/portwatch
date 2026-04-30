package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80, State: "LISTEN"},
		{Protocol: "tcp", Address: "127.0.0.1", Port: 443, State: "LISTEN"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	ports := samplePorts()
	if err := store.Save(ports); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap == nil {
		t.Fatal("Load() returned nil snapshot, want non-nil")
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("got %d ports, want %d", len(snap.Ports), len(ports))
	}
	if snap.Timestamp.IsZero() {
		t.Error("snapshot timestamp should not be zero")
	}
	if time.Since(snap.Timestamp) > 5*time.Second {
		t.Error("snapshot timestamp seems too old")
	}
}

func TestLoadNoFile(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() on missing file returned error: %v", err)
	}
	if snap != nil {
		t.Errorf("Load() = %v, want nil when no snapshot exists", snap)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/portwatch"
	store := snapshot.NewStore(dir)

	if err := store.Save(samplePorts()); err != nil {
		t.Fatalf("Save() failed to create nested dir: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSaveOverwritesPreviousSnapshot(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	first := samplePorts()
	if err := store.Save(first); err != nil {
		t.Fatalf("first Save() error: %v", err)
	}

	second := []scanner.Port{{Protocol: "udp", Address: "0.0.0.0", Port: 53, State: "LISTEN"}}
	if err := store.Save(second); err != nil {
		t.Fatalf("second Save() error: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(snap.Ports) != 1 {
		t.Errorf("got %d ports after overwrite, want 1", len(snap.Ports))
	}
}
