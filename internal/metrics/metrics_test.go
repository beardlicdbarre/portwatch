package metrics

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNewCollectorStartedAt(t *testing.T) {
	before := time.Now()
	c := New()
	after := time.Now()
	s := c.Snapshot()
	if s.StartedAt.Before(before) || s.StartedAt.After(after) {
		t.Errorf("StartedAt %v outside [%v, %v]", s.StartedAt, before, after)
	}
}

func TestRecordScanIncrementsCounter(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	if got := c.Snapshot().ScansTotal; got != 2 {
		t.Errorf("ScansTotal = %d, want 2", got)
	}
}

func TestRecordChangeIncrementsCounter(t *testing.T) {
	c := New()
	c.RecordChange()
	if got := c.Snapshot().ChangesTotal; got != 1 {
		t.Errorf("ChangesTotal = %d, want 1", got)
	}
}

func TestRecordAlertIncrementsCounter(t *testing.T) {
	c := New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	if got := c.Snapshot().AlertsTotal; got != 3 {
		t.Errorf("AlertsTotal = %d, want 3", got)
	}
}

func TestLastScanAtUpdated(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	c := New()
	c.now = fixedNow(base)
	c.RecordScan()
	if got := c.Snapshot().LastScanAt; !got.Equal(base) {
		t.Errorf("LastScanAt = %v, want %v", got, base)
	}
}

func TestLastChangeAtUpdated(t *testing.T) {
	base := time.Date(2024, 6, 15, 8, 30, 0, 0, time.UTC)
	c := New()
	c.now = fixedNow(base)
	c.RecordChange()
	if got := c.Snapshot().LastChangeAt; !got.Equal(base) {
		t.Errorf("LastChangeAt = %v, want %v", got, base)
	}
}

func TestUptimeSeconds(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c := &Collector{
		startedAt: start,
		now:       fixedNow(start.Add(5 * time.Second)),
	}
	s := c.Snapshot()
	if s.UptimeSeconds != 5.0 {
		t.Errorf("UptimeSeconds = %f, want 5.0", s.UptimeSeconds)
	}
}

func TestSnapshotIsImmutable(t *testing.T) {
	c := New()
	s1 := c.Snapshot()
	c.RecordScan()
	s2 := c.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Error("expected snapshots to differ after RecordScan")
	}
}
