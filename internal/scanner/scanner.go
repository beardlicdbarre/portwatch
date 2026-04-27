package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Port represents an open port with its protocol and process info.
type Port struct {
	Protocol string
	Address  string
	Port     int
	State    string
}

// String returns a human-readable representation of a Port.
func (p Port) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Key returns a unique identifier for the port.
func (p Port) Key() string {
	return fmt.Sprintf("%s/%s:%d", p.Protocol, p.Address, p.Port)
}

// Scanner defines the interface for port scanning backends.
type Scanner interface {
	Scan() ([]Port, error)
}

// ParseAddress splits a combined host:port string into its components.
func ParseAddress(addr string) (host string, port int, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		// Try treating the whole string as a port number
		portStr = addr
		host = ""
	}
	portStr = strings.TrimSpace(portStr)
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port %q: %w", portStr, err)
	}
	if port < 1 || port > 65535 {
		return "", 0, fmt.Errorf("port %d out of valid range (1-65535)", port)
	}
	return host, port, nil
}

// Diff computes added and removed ports between two snapshots.
func Diff(previous, current []Port) (added, removed []Port) {
	prevMap := make(map[string]Port, len(previous))
	for _, p := range previous {
		prevMap[p.Key()] = p
	}
	currMap := make(map[string]Port, len(current))
	for _, p := range current {
		currMap[p.Key()] = p
	}
	for key, p := range currMap {
		if _, exists := prevMap[key]; !exists {
			added = append(added, p)
		}
	}
	for key, p := range prevMap {
		if _, exists := currMap[key]; !exists {
			removed = append(removed, p)
		}
	}
	return added, removed
}
