package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon runs the port monitoring loop.
type Daemon struct {
	cfg     *config.Config
	scanner scanner.Scanner
	notifier *alert.Notifier
}

// New creates a new Daemon with the given configuration.
func New(cfg *config.Config, sc scanner.Scanner, n *alert.Notifier) *Daemon {
	return &Daemon{
		cfg:      cfg,
		scanner:  sc,
		notifier: n,
	}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
// It performs an initial scan to establish a baseline, then scans at
// the interval defined in cfg. Any changes detected between scans are
// forwarded to the notifier. Scan or notify errors are logged but do
// not stop the loop; only context cancellation causes Run to return.
func (d *Daemon) Run(ctx context.Context) error {
	prev, err := d.scanner.Scan()
	if err != nil {
		return err
	}
	log.Printf("portwatch: initial scan found %d open ports", len(prev))

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch: shutting down")
			return ctx.Err()
		case <-ticker.C:
			d.tick(ctx, &prev)
		}
	}
}

// tick performs a single scan cycle: scans ports, diffs against prev,
// notifies on changes, and updates prev to the current snapshot.
func (d *Daemon) tick(ctx context.Context, prev *scanner.PortSet) {
	curr, err := d.scanner.Scan()
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return
	}
	diff := scanner.Diff(*prev, curr)
	if len(diff.Opened)+len(diff.Closed) > 0 {
		if err := d.notifier.Notify(diff); err != nil {
			log.Printf("portwatch: notify error: %v", err)
		}
	}
	*prev = curr
}
