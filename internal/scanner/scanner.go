package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Port represents a single listening network port.
type Port struct {
	Protocol string
	Address  string
	PortNum  int
}

// PortKey returns a unique string key for a Port.
func (p Port) PortKey() string {
	return fmt.Sprintf("%s:%s:%d", p.Protocol, p.Address, p.PortNum)
}

// String returns a human-readable representation of a Port.
func (p Port) String() string {
	return fmt.Sprintf("%s/%s:%d", p.Protocol, p.Address, p.PortNum)
}

// DiffResult holds the changes between two port snapshots.
type DiffResult struct {
	Added   []Port
	Removed []Port
}

// HasChanges returns true if there are any added or removed ports.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Scanner is the interface implemented by port scanning backends.
type Scanner interface {
	Scan() ([]Port, error)
}

// ParseAddress parses a combined "host:port" string into its components.
func ParseAddress(raw string) (host string, port int, err error) {
	raw = strings.TrimSpace(raw)
	h, p, err := net.SplitHostPort(raw)
	if err != nil {
		return "", 0, fmt.Errorf("invalid address %q: %w", raw, err)
	}
	portNum, err := strconv.Atoi(p)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port in %q: %w", raw, err)
	}
	return h, portNum, nil
}

// Diff compares two port snapshots and returns what was added or removed.
func Diff(previous, current []Port) DiffResult {
	prev := make(map[string]Port, len(previous))
	for _, p := range previous {
		prev[p.PortKey()] = p
	}

	curr := make(map[string]Port, len(current))
	for _, p := range current {
		curr[p.PortKey()] = p
	}

	var result DiffResult
	for k, p := range curr {
		if _, exists := prev[k]; !exists {
			result.Added = append(result.Added, p)
		}
	}
	for k, p := range prev {
		if _, exists := curr[k]; !exists {
			result.Removed = append(result.Removed, p)
		}
	}
	return result
}
