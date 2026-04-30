package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same port key within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, the key's timestamp is updated.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if last, ok := l.seen[key]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}
	l.seen[key] = now
	return true
}

// Reset clears the recorded timestamp for a key, allowing it immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.seen, key)
}

// Purge removes all entries older than the cooldown window to keep memory bounded.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := l.now().Add(-l.cooldown)
	for k, t := range l.seen {
		if t.Before(cutoff) {
			delete(l.seen, k)
		}
	}
}
