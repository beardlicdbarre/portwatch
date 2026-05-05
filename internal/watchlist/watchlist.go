// Package watchlist manages a set of ports that are expected to always
// be open. Any port in the watchlist that disappears triggers a high-
// priority alert regardless of the normal filter or rate-limit rules.
package watchlist

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Entry is a single watchlist record: a protocol + port number pair.
type Entry struct {
	Proto string
	Port  uint16
}

// String returns a human-readable representation such as "tcp:443".
func (e Entry) String() string {
	return fmt.Sprintf("%s:%d", strings.ToLower(e.Proto), e.Port)
}

// Watchlist holds the set of ports that must remain open.
type Watchlist struct {
	entries map[string]Entry
}

// New creates a Watchlist from a slice of Entry values.
func New(entries []Entry) *Watchlist {
	wl := &Watchlist{entries: make(map[string]Entry, len(entries))}
	for _, e := range entries {
		wl.entries[e.String()] = e
	}
	return wl
}

// Missing returns all watchlist entries whose port is not present in ports.
func (w *Watchlist) Missing(ports []scanner.Port) []Entry {
	present := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		key := fmt.Sprintf("%s:%d", strings.ToLower(p.Proto), p.Port)
		present[key] = struct{}{}
	}

	var missing []Entry
	for key, e := range w.entries {
		if _, ok := present[key]; !ok {
			missing = append(missing, e)
		}
	}
	return missing
}

// Contains reports whether the given port is on the watchlist.
func (w *Watchlist) Contains(p scanner.Port) bool {
	key := fmt.Sprintf("%s:%d", strings.ToLower(p.Proto), p.Port)
	_, ok := w.entries[key]
	return ok
}

// Len returns the number of entries in the watchlist.
func (w *Watchlist) Len() int { return len(w.entries) }
