// Package suppress provides a mechanism to temporarily silence alerts
// for specific ports, preventing alert fatigue during known maintenance windows.
package suppress

import (
	"sync"
	"time"
)

// Key identifies a suppressed port by its address and protocol.
type Key struct {
	Address  string
	Protocol string
}

// entry holds the expiry time for a suppression rule.
type entry struct {
	expiry time.Time
}

// List manages a set of time-bounded suppression rules.
type List struct {
	mu      sync.Mutex
	rules   map[Key]entry
	clock   func() time.Time
}

// New returns a new List using the real wall clock.
func New() *List {
	return &List{
		rules: make(map[Key]entry),
		clock: time.Now,
	}
}

// newWithClock returns a List with an injectable clock for testing.
func newWithClock(clock func() time.Time) *List {
	return &List{
		rules: make(map[Key]entry),
		clock: clock,
	}
}

// Add suppresses alerts for the given key for the specified duration.
// Calling Add on an already-suppressed key extends or replaces the expiry.
func (l *List) Add(k Key, duration time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rules[k] = entry{expiry: l.clock().Add(duration)}
}

// Remove lifts a suppression rule immediately.
func (l *List) Remove(k Key) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.rules, k)
}

// IsSuppressed reports whether the given key is currently suppressed.
// Expired rules are lazily pruned on access.
func (l *List) IsSuppressed(k Key) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.rules[k]
	if !ok {
		return false
	}
	if l.clock().After(e.expiry) {
		delete(l.rules, k)
		return false
	}
	return true
}

// Len returns the number of active (non-expired) suppression rules.
func (l *List) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock()
	count := 0
	for k, e := range l.rules {
		if now.After(e.expiry) {
			delete(l.rules, k)
		} else {
			count++
		}
	}
	return count
}
