package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a saved port state with metadata.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Ports     []scanner.Port     `json:"ports"`
}

// Store persists and loads port snapshots to/from disk.
type Store struct {
	dir string
}

// NewStore creates a Store that reads/writes snapshots under dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes the current port list as a snapshot to disk.
func (s *Store) Save(ports []scanner.Port) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir %s: %w", s.dir, err)
	}

	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	path := s.filePath()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads the latest snapshot from disk. Returns nil, nil if no snapshot exists.
func (s *Store) Load() (*Snapshot, error) {
	path := s.filePath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

func (s *Store) filePath() string {
	return filepath.Join(s.dir, "portwatch.snapshot.json")
}
