package filter

import (
	"net"

	"github.com/user/portwatch/internal/scanner"
)

// Rule defines a filter rule that can exclude ports from alerting.
type Rule struct {
	Protocol string // "tcp", "udp", or "" for any
	Port     uint16 // 0 means any port
	Host     string // empty means any host
}

// Filter decides which port changes should trigger alerts.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Allow returns true if the port should be included in alerting
// (i.e. it is NOT matched by any exclusion rule).
func (f *Filter) Allow(p scanner.Port) bool {
	for _, r := range f.rules {
		if f.matches(r, p) {
			return false
		}
	}
	return true
}

// FilterDiff removes ports from added/removed slices that match exclusion rules.
func (f *Filter) FilterDiff(added, removed []scanner.Port) ([]scanner.Port, []scanner.Port) {
	return filterSlice(f, added), filterSlice(f, removed)
}

func filterSlice(f *Filter, ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		if f.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}

func (f *Filter) matches(r Rule, p scanner.Port) bool {
	if r.Protocol != "" && r.Protocol != p.Protocol {
		return false
	}
	if r.Port != 0 && r.Port != p.Port {
		return false
	}
	if r.Host != "" {
		rIP := net.ParseIP(r.Host)
		pIP := net.ParseIP(p.Host)
		if rIP == nil || pIP == nil {
			if r.Host != p.Host {
				return false
			}
		} else if !rIP.Equal(pIP) {
			return false
		}
	}
	return true
}
