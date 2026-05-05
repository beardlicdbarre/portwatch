package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func sampleSnapshot() Snapshot {
	return Snapshot{
		ScansTotal:    42,
		ChangesTotal:  7,
		AlertsTotal:   3,
		StartedAt:     time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC),
		LastScanAt:    time.Date(2024, 3, 10, 9, 5, 0, 0, time.UTC),
		LastChangeAt:  time.Date(2024, 3, 10, 9, 4, 0, 0, time.UTC),
		UptimeSeconds: 300,
	}
}

func TestReportTextContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, FormatText)
	if err := r.Report(sampleSnapshot()); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"scans_total", "42", "changes_total", "7", "alerts_total", "3", "uptime_seconds", "300"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestReportJSONIsValid(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, FormatJSON)
	if err := r.Report(sampleSnapshot()); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	var s Snapshot
	if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if s.ScansTotal != 42 {
		t.Errorf("ScansTotal = %d, want 42", s.ScansTotal)
	}
}

func TestReportTextOmitsZeroLastScanAt(t *testing.T) {
	snap := sampleSnapshot()
	snap.LastScanAt = time.Time{}
	var buf bytes.Buffer
	r := NewReporter(&buf, FormatText)
	_ = r.Report(snap)
	if strings.Contains(buf.String(), "last_scan_at") {
		t.Error("expected last_scan_at to be omitted when zero")
	}
}

func TestReportDefaultFormatIsText(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf, "")
	if err := r.Report(sampleSnapshot()); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	if strings.HasPrefix(strings.TrimSpace(buf.String()), "{") {
		t.Error("expected text output, got JSON")
	}
}
