package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event represents a single port change event recorded in history.
type Event struct {
	Timestamp time.Time    `json:"timestamp"`
	Kind      string       `json:"kind"` // "added" or "removed"
	Port      scanner.Port `json:"port"`
}

// Store persists port change events to a JSON-lines log file.
type Store struct {
	mu   sync.Mutex
	path string
}

// NewStore creates a new history Store that writes to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Record appends a new event to the history log.
func (s *Store) Record(kind string, port scanner.Port, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	ev := Event{Timestamp: ts, Kind: kind, Port: port}
	enc := json.NewEncoder(f)
	return enc.Encode(ev)
}

// Load reads all events from the history log.
func (s *Store) Load() ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []Event
	dec := json.NewDecoder(f)
	for dec.More() {
		var ev Event
		if err := dec.Decode(&ev); err != nil {
			return events, err
		}
		events = append(events, ev)
	}
	return events, nil
}
