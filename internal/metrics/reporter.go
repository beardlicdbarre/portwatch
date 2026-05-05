package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Format controls how a Snapshot is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes metric snapshots to an io.Writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// NewReporter returns a Reporter that writes to w using the given format.
func NewReporter(w io.Writer, format Format) *Reporter {
	return &Reporter{w: w, format: format}
}

// Report renders the snapshot to the underlying writer.
func (r *Reporter) Report(s Snapshot) error {
	switch r.format {
	case FormatJSON:
		return r.reportJSON(s)
	default:
		return r.reportText(s)
	}
}

func (r *Reporter) reportText(s Snapshot) error {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "started_at\t%s\n", s.StartedAt.Format(time.RFC3339))
	fmt.Fprintf(tw, "uptime_seconds\t%.0f\n", s.UptimeSeconds)
	fmt.Fprintf(tw, "scans_total\t%d\n", s.ScansTotal)
	fmt.Fprintf(tw, "changes_total\t%d\n", s.ChangesTotal)
	fmt.Fprintf(tw, "alerts_total\t%d\n", s.AlertsTotal)
	if !s.LastScanAt.IsZero() {
		fmt.Fprintf(tw, "last_scan_at\t%s\n", s.LastScanAt.Format(time.RFC3339))
	}
	if !s.LastChangeAt.IsZero() {
		fmt.Fprintf(tw, "last_change_at\t%s\n", s.LastChangeAt.Format(time.RFC3339))
	}
	return tw.Flush()
}

func (r *Reporter) reportJSON(s Snapshot) error {
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
