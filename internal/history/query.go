package history

import (
	"time"
)

// Filter holds optional criteria for querying history events.
type Filter struct {
	Kind  string    // "added", "removed", or "" for all
	Since time.Time // zero means no lower bound
	Proto string    // "tcp", "udp", or "" for all
}

// Query returns events from the store that match the given filter.
func (s *Store) Query(f Filter) ([]Event, error) {
	all, err := s.Load()
	if err != nil {
		return nil, err
	}

	var out []Event
	for _, ev := range all {
		if f.Kind != "" && ev.Kind != f.Kind {
			continue
		}
		if !f.Since.IsZero() && ev.Timestamp.Before(f.Since) {
			continue
		}
		if f.Proto != "" && ev.Port.Proto != f.Proto {
			continue
		}
		out = append(out, ev)
	}
	return out, nil
}
