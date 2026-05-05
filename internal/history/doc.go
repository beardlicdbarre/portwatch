// Package history provides persistent event logging for portwatch.
//
// It records port change events (additions and removals) to a JSON-lines
// file on disk, and exposes a query interface for filtering events by
// kind, protocol, or time range.
//
// Usage:
//
//	store := history.NewStore("/var/lib/portwatch/history.jsonl")
//
//	// Record an event
//	store.Record("added", port, time.Now())
//
//	// Query recent TCP additions
//	events, err := store.Query(history.Filter{
//		Kind:  "added",
//		Proto: "tcp",
//		Since: time.Now().Add(-24 * time.Hour),
//	})
package history
