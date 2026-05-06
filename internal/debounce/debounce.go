// Package debounce provides a mechanism to suppress rapid repeated events
// for the same port key, emitting only after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Clock allows injecting a time source for testing.
type Clock func() time.Time

// Debouncer holds pending events and releases them only after the quiet
// window has passed without a repeat for that key.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	clock   Clock
	pending map[string]time.Time
}

// New returns a Debouncer that suppresses repeated events within window.
func New(window time.Duration) *Debouncer {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, clock Clock) *Debouncer {
	return &Debouncer{
		window:  window,
		clock:   clock,
		pending: make(map[string]time.Time),
	}
}

// Seen records that an event for key was observed at the current time.
// Call this every time the event fires.
func (d *Debouncer) Seen(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pending[key] = d.clock()
}

// Ready reports whether the event for key has been quiet for at least the
// configured window, meaning it is safe to emit. If Ready returns true the
// pending entry is cleared so the next Seen starts a fresh window.
func (d *Debouncer) Ready(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	last, ok := d.pending[key]
	if !ok {
		return false
	}
	if d.clock().Sub(last) >= d.window {
		delete(d.pending, key)
		return true
	}
	return false
}

// Forget removes any pending state for key without emitting.
func (d *Debouncer) Forget(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.pending, key)
}

// PendingCount returns the number of keys currently waiting in the window.
func (d *Debouncer) PendingCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}
