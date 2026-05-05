package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Proto: proto, Port: port}
}

func TestContainsMatchingEntry(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{
		{Proto: "tcp", Port: 443},
	})
	if !wl.Contains(makePort("tcp", 443)) {
		t.Error("expected watchlist to contain tcp:443")
	}
}

func TestContainsCaseInsensitive(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{{Proto: "TCP", Port: 80}})
	if !wl.Contains(makePort("tcp", 80)) {
		t.Error("expected case-insensitive match for tcp:80")
	}
}

func TestContainsMiss(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{{Proto: "tcp", Port: 443}})
	if wl.Contains(makePort("tcp", 80)) {
		t.Error("tcp:80 should not be in watchlist")
	}
}

func TestMissingAllPresent(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{
		{Proto: "tcp", Port: 443},
		{Proto: "tcp", Port: 80},
	})
	ports := []scanner.Port{makePort("tcp", 443), makePort("tcp", 80)}
	if got := wl.Missing(ports); len(got) != 0 {
		t.Errorf("expected no missing entries, got %v", got)
	}
}

func TestMissingDetectsAbsent(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{
		{Proto: "tcp", Port: 443},
		{Proto: "tcp", Port: 22},
	})
	ports := []scanner.Port{makePort("tcp", 443)}
	missing := wl.Missing(ports)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing entry, got %d", len(missing))
	}
	if missing[0].Port != 22 {
		t.Errorf("expected missing port 22, got %d", missing[0].Port)
	}
}

func TestMissingEmptyPorts(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{
		{Proto: "udp", Port: 53},
	})
	missing := wl.Missing(nil)
	if len(missing) != 1 {
		t.Errorf("expected 1 missing entry, got %d", len(missing))
	}
}

func TestLen(t *testing.T) {
	wl := watchlist.New([]watchlist.Entry{
		{Proto: "tcp", Port: 80},
		{Proto: "tcp", Port: 443},
		{Proto: "udp", Port: 53},
	})
	if wl.Len() != 3 {
		t.Errorf("expected Len 3, got %d", wl.Len())
	}
}

func TestEntryString(t *testing.T) {
	e := watchlist.Entry{Proto: "TCP", Port: 8080}
	if got := e.String(); got != "tcp:8080" {
		t.Errorf("expected \"tcp:8080\", got %q", got)
	}
}
