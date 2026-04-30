package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon orchestrates periodic port scanning and alerting.
type Daemon struct {
	cfg      *config.Config
	scanner  scanner.Scanner
	store    *snapshot.Store
	notifier *alert.Notifier
	filter   *filter.Filter
	limiter  *ratelimit.Limiter
}

// New creates a Daemon wired with all dependencies.
func New(cfg *config.Config, sc scanner.Scanner, store *snapshot.Store, n *alert.Notifier, f *filter.Filter) *Daemon {
	return &Daemon{
		cfg:     cfg,
		scanner: sc,
		store:   store,
		notifier: n,
		filter:  f,
		limiter: ratelimit.New(cfg.AlertCooldown),
	}
}

// Run starts the scan loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	prev, err := d.store.Load()
	if err != nil {
		log.Printf("portwatch: no previous snapshot, starting fresh: %v", err)
		prev = nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			prev = d.tick(prev)
		}
	}
}

func (d *Daemon) tick(prev []scanner.Port) []scanner.Port {
	current, err := d.scanner.Scan()
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return prev
	}

	current = d.filter.Apply(current)
	diff := scanner.Diff(prev, current)

	for _, p := range diff.Added {
		if d.limiter.Allow(p.Key()) {
			d.notifier.Notify(diff)
			break
		}
	}
	for _, p := range diff.Removed {
		if d.limiter.Allow(p.Key()) {
			d.notifier.Notify(diff)
			break
		}
	}

	if err := d.store.Save(current); err != nil {
		log.Printf("portwatch: snapshot save error: %v", err)
	}
	return current
}
