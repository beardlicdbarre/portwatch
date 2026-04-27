package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.Port
}

// Notifier sends alerts to a configured destination.
type Notifier struct {
	Writer io.Writer
}

// NewNotifier creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{Writer: w}
}

// Notify formats and writes an alert for each port change in the diff.
func (n *Notifier) Notify(diff scanner.DiffResult) error {
	for _, p := range diff.Added {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   "new port opened",
			Port:      p,
		}
		if err := n.write(a); err != nil {
			return err
		}
	}
	for _, p := range diff.Removed {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   "port closed",
			Port:      p,
		}
		if err := n.write(a); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) write(a Alert) error {
	_, err := fmt.Fprintf(
		n.Writer,
		"[%s] %s %-6s %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port.Protocol,
		a.Port.String(),
	)
	return err
}
