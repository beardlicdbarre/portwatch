package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func port(proto, host string, p uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Host: host, Port: p}
}

func TestAllowNoRules(t *testing.T) {
	f := filter.New(nil)
	if !f.Allow(port("tcp", "0.0.0.0", 8080)) {
		t.Error("expected port to be allowed with no rules")
	}
}

func TestAllowMatchingPortBlocked(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Protocol: "tcp", Port: 22},
	})
	if f.Allow(port("tcp", "0.0.0.0", 22)) {
		t.Error("expected port 22/tcp to be blocked")
	}
}

func TestAllowNonMatchingPortPasses(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Protocol: "tcp", Port: 22},
	})
	if !f.Allow(port("tcp", "0.0.0.0", 80)) {
		t.Error("expected port 80/tcp to be allowed")
	}
}

func TestAllowProtocolMismatch(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Protocol: "tcp", Port: 53},
	})
	// UDP 53 should still be allowed
	if !f.Allow(port("udp", "0.0.0.0", 53)) {
		t.Error("expected udp/53 to be allowed when rule only covers tcp/53")
	}
}

func TestAllowHostFilter(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Host: "127.0.0.1", Port: 9000},
	})
	if f.Allow(port("tcp", "127.0.0.1", 9000)) {
		t.Error("expected 127.0.0.1:9000 to be blocked")
	}
	if !f.Allow(port("tcp", "0.0.0.0", 9000)) {
		t.Error("expected 0.0.0.0:9000 to be allowed")
	}
}

func TestFilterDiff(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Protocol: "tcp", Port: 22},
	})
	added := []scanner.Port{
		port("tcp", "0.0.0.0", 22),
		port("tcp", "0.0.0.0", 8080),
	}
	removed := []scanner.Port{
		port("tcp", "0.0.0.0", 22),
	}

	filteredAdded, filteredRemoved := f.FilterDiff(added, removed)

	if len(filteredAdded) != 1 || filteredAdded[0].Port != 8080 {
		t.Errorf("expected only port 8080 in added, got %v", filteredAdded)
	}
	if len(filteredRemoved) != 0 {
		t.Errorf("expected empty removed, got %v", filteredRemoved)
	}
}
