package suppress

import (
	"testing"
	"time"
)

var (
	now   = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	tcpKey = Key{Address: "0.0.0.0:8080", Protocol: "tcp"}
	udpKey = Key{Address: "0.0.0.0:53", Protocol: "udp"}
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestIsSuppressedAfterAdd(t *testing.T) {
	l := newWithClock(fixedClock(now))
	l.Add(tcpKey, 10*time.Minute)
	if !l.IsSuppressed(tcpKey) {
		t.Fatal("expected key to be suppressed")
	}
}

func TestNotSuppressedWithoutAdd(t *testing.T) {
	l := newWithClock(fixedClock(now))
	if l.IsSuppressed(tcpKey) {
		t.Fatal("expected key to not be suppressed")
	}
}

func TestExpiryLiftsSuppressionAutomatically(t *testing.T) {
	current := now
	l := newWithClock(func() time.Time { return current })
	l.Add(tcpKey, 5*time.Minute)

	// advance past expiry
	current = now.Add(6 * time.Minute)
	if l.IsSuppressed(tcpKey) {
		t.Fatal("expected suppression to have expired")
	}
}

func TestRemoveLiftsSuppression(t *testing.T) {
	l := newWithClock(fixedClock(now))
	l.Add(tcpKey, 1*time.Hour)
	l.Remove(tcpKey)
	if l.IsSuppressed(tcpKey) {
		t.Fatal("expected suppression to be removed")
	}
}

func TestAddExtendsExpiry(t *testing.T) {
	current := now
	l := newWithClock(func() time.Time { return current })
	l.Add(tcpKey, 2*time.Minute)

	// re-add with longer duration
	l.Add(tcpKey, 10*time.Minute)

	// advance past original expiry but within new one
	current = now.Add(3 * time.Minute)
	if !l.IsSuppressed(tcpKey) {
		t.Fatal("expected suppression to still be active after extension")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	l := newWithClock(fixedClock(now))
	l.Add(tcpKey, 10*time.Minute)

	if l.IsSuppressed(udpKey) {
		t.Fatal("udpKey should not be suppressed")
	}
	if !l.IsSuppressed(tcpKey) {
		t.Fatal("tcpKey should be suppressed")
	}
}

func TestLenCountsActiveRules(t *testing.T) {
	current := now
	l := newWithClock(func() time.Time { return current })
	l.Add(tcpKey, 5*time.Minute)
	l.Add(udpKey, 1*time.Minute)

	if got := l.Len(); got != 2 {
		t.Fatalf("expected 2 active rules, got %d", got)
	}

	// expire udpKey
	current = now.Add(2 * time.Minute)
	if got := l.Len(); got != 1 {
		t.Fatalf("expected 1 active rule after expiry, got %d", got)
	}
}
