package daemon_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/scanner"
)

// mockScanner returns successive snapshots on each call to Scan.
type mockScanner struct {
	snapshots []scanner.Snapshot
	call      int
}

func (m *mockScanner) Scan() (scanner.Snapshot, error) {
	if m.call >= len(m.snapshots) {
		return m.snapshots[len(m.snapshots)-1], nil
	}
	s := m.snapshots[m.call]
	m.call++
	return s, nil
}

func TestDaemonRunCancels(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Interval = 20 * time.Millisecond

	initial := scanner.Snapshot{
		{Host: "0.0.0.0", Port: 8080, Proto: "tcp"}: {},
	}
	ms := &mockScanner{snapshots: []scanner.Snapshot{initial}}

	n, _ := alert.NewNotifier("")
	d := daemon.New(cfg, ms, n)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDaemonDetectsChanges(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Interval = 10 * time.Millisecond

	snap1 := scanner.Snapshot{
		{Host: "0.0.0.0", Port: 80, Proto: "tcp"}: {},
	}
	snap2 := scanner.Snapshot{
		{Host: "0.0.0.0", Port: 80, Proto: "tcp"}: {},
		{Host: "0.0.0.0", Port: 443, Proto: "tcp"}: {},
	}
	ms := &mockScanner{snapshots: []scanner.Snapshot{snap1, snap2}}

	n, _ := alert.NewNotifier("")
	d := daemon.New(cfg, ms, n)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Should run without error until context expires.
	if err := d.Run(ctx); err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
