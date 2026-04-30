package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls the output format for alerts.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter converts port diff events into human-readable or structured output.
type Formatter struct {
	format    Format
	timestamp bool
}

// New creates a Formatter with the given format and timestamp option.
func New(format Format, timestamp bool) *Formatter {
	return &Formatter{format: format, timestamp: timestamp}
}

// FormatAdded returns a formatted string for a newly detected port.
func (f *Formatter) FormatAdded(p scanner.Port) string {
	return f.formatEvent("ADDED", p)
}

// FormatRemoved returns a formatted string for a port that disappeared.
func (f *Formatter) FormatRemoved(p scanner.Port) string {
	return f.formatEvent("REMOVED", p)
}

func (f *Formatter) formatEvent(event string, p scanner.Port) string {
	switch f.format {
	case FormatJSON:
		return f.jsonLine(event, p)
	default:
		return f.textLine(event, p)
	}
}

func (f *Formatter) textLine(event string, p scanner.Port) string {
	parts := []string{}
	if f.timestamp {
		parts = append(parts, time.Now().UTC().Format(time.RFC3339))
	}
	parts = append(parts, fmt.Sprintf("[%s]", event))
	parts = append(parts, p.String())
	return strings.Join(parts, " ")
}

func (f *Formatter) jsonLine(event string, p scanner.Port) string {
	ts := ""
	if f.timestamp {
		ts = fmt.Sprintf(`"timestamp":%q,`, time.Now().UTC().Format(time.RFC3339))
	}
	return fmt.Sprintf(`{%s"event":%q,"proto":%q,"addr":%q,"port":%d}`,
		ts, event, p.Proto, p.Addr, p.Port)
}
