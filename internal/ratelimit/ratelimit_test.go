package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowFirstTime(t *testing.T) {
	l := New(5 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowBlockedWithinCooldown(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("tcp:8080")

	l.now = fixedClock(base.Add(5 * time.Second))
	if l.Allow("tcp:8080") {
		t.Fatal("expected call within cooldown to be blocked")
	}
}

func TestAllowAfterCooldownExpires(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("tcp:8080")

	l.now = fixedClock(base.Add(11 * time.Second))
	if !l.Allow("tcp:8080") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllowDifferentKeysIndependent(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("tcp:8080")

	if !l.Allow("udp:9090") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("tcp:8080")
	l.Reset("tcp:8080")

	if !l.Allow("tcp:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestPurgeRemovesStaleEntries(t *testing.T) {
	base := time.Now()
	l := New(5 * time.Second)
	l.now = fixedClock(base)
	l.Allow("tcp:8080")

	l.now = fixedClock(base.Add(10 * time.Second))
	l.Purge()

	if len(l.seen) != 0 {
		t.Fatalf("expected seen map to be empty after purge, got %d entries", len(l.seen))
	}
}

func TestPurgeKeepsFreshEntries(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)
	l.Allow("tcp:8080")

	l.now = fixedClock(base.Add(5 * time.Second))
	l.Purge()

	if len(l.seen) != 1 {
		t.Fatalf("expected 1 entry to remain after purge, got %d", len(l.seen))
	}
}
