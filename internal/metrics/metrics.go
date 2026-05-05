// Package metrics tracks runtime statistics for the portwatch daemon.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of daemon metrics.
type Snapshot struct {
	ScansTotal    int64
	ChangesTotal  int64
	AlertsTotal   int64
	LastScanAt    time.Time
	LastChangeAt  time.Time
	UptimeSeconds float64
	StartedAt     time.Time
}

// Collector accumulates runtime metrics in a thread-safe manner.
type Collector struct {
	mu           sync.RWMutex
	scansTotal   int64
	changesTotal int64
	alertsTotal  int64
	lastScanAt   time.Time
	lastChangeAt time.Time
	startedAt    time.Time
	now          func() time.Time
}

// New returns a new Collector initialised with the current time.
func New() *Collector {
	n := time.Now()
	return &Collector{
		startedAt: n,
		now:       time.Now,
	}
}

// RecordScan increments the scan counter and records the scan timestamp.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = c.now()
}

// RecordChange increments the change counter and records the change timestamp.
func (c *Collector) RecordChange() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changesTotal++
	c.lastChangeAt = c.now()
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns an immutable copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		ChangesTotal:  c.changesTotal,
		AlertsTotal:   c.alertsTotal,
		LastScanAt:    c.lastScanAt,
		LastChangeAt:  c.lastChangeAt,
		UptimeSeconds: c.now().Sub(c.startedAt).Seconds(),
		StartedAt:     c.startedAt,
	}
}
