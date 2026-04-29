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
			curr, err := d.scanner.Scan()
			if err != nil {
				log.Printf("portwatch: scan error: %v", err)
				continue
			}
			diff := scanner.Diff(prev, curr)
			if err := d.notifier.Notify(diff); err != nil {
				log.Printf("portwatch: notify error: %v", err)
			}
			prev = curr
		}
	}
}
