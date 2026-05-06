package debounce

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type fixedClock struct {
	now time.Time
}

func (f *fixedClock) tick(d time.Duration) { f.now = f.now.Add(d) }
func (f *fixedClock) get() time.Time       { return f.now }

func TestReadyReturnsFalseWithoutSeen(t *testing.T) {
	d := New(time.Second)
	if d.Ready("tcp:8080") {
		t.Fatal("expected Ready to be false for unseen key")
	}
}

func TestReadyReturnsFalseWithinWindow(t *testing.T) {
	clk := &fixedClock{now: epoch}
	d := newWithClock(2*time.Second, clk.get)

	d.Seen("tcp:8080")
	clk.tick(1 * time.Second)

	if d.Ready("tcp:8080") {
		t.Fatal("expected Ready to be false within debounce window")
	}
}

func TestReadyReturnsTrueAfterWindow(t *testing.T) {
	clk := &fixedClock{now: epoch}
	d := newWithClock(2*time.Second, clk.get)

	d.Seen("tcp:8080")
	clk.tick(2 * time.Second)

	if !d.Ready("tcp:8080") {
		t.Fatal("expected Ready to be true after debounce window")
	}
}

func TestReadyClearsPendingEntry(t *testing.T) {
	clk := &fixedClock{now: epoch}
	d := newWithClock(time.Second, clk.get)

	d.Seen("tcp:9090")
	clk.tick(2 * time.Second)
	d.Ready("tcp:9090") // consume

	if d.PendingCount() != 0 {
		t.Fatalf("expected 0 pending after Ready, got %d", d.PendingCount())
	}
}

func TestSeenResetsWindow(t *testing.T) {
	clk := &fixedClock{now: epoch}
	d := newWithClock(2*time.Second, clk.get)

	d.Seen("tcp:443")
	clk.tick(1 * time.Second)
	d.Seen("tcp:443") // reset window
	clk.tick(1 * time.Second)

	if d.Ready("tcp:443") {
		t.Fatal("expected Ready false — window was reset by second Seen")
	}

	clk.tick(1 * time.Second) // now 2 s since last Seen
	if !d.Ready("tcp:443") {
		t.Fatal("expected Ready true after window elapsed from reset")
	}
}

func TestForgetRemovesPendingKey(t *testing.T) {
	clk := &fixedClock{now: epoch}
	d := newWithClock(time.Second, clk.get)

	d.Seen("udp:53")
	d.Forget("udp:53")

	if d.PendingCount() != 0 {
		t.Fatalf("expected 0 pending after Forget, got %d", d.PendingCount())
	}
	if d.Ready("udp:53") {
		t.Fatal("expected Ready false after Forget")
	}
}

func TestPendingCountTracksMultipleKeys(t *testing.T) {
	d := New(time.Minute)
	d.Seen("tcp:80")
	d.Seen("tcp:443")
	d.Seen("udp:53")

	if got := d.PendingCount(); got != 3 {
		t.Fatalf("expected 3 pending, got %d", got)
	}
}
