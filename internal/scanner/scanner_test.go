package scanner

import (
	"testing"
)

func TestPortKey(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "0.0.0.0", Port: 8080}
	expected := "tcp/0.0.0.0:8080"
	if p.Key() != expected {
		t.Errorf("Key() = %q, want %q", p.Key(), expected)
	}
}

func TestPortString(t *testing.T) {
	p := Port{Protocol: "tcp", Address: "127.0.0.1", Port: 443}
	expected := "127.0.0.1:443 (tcp)"
	if p.String() != expected {
		t.Errorf("String() = %q, want %q", p.String(), expected)
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		input   string
		host    string
		port    int
		wantErr bool
	}{
		{"0.0.0.0:80", "0.0.0.0", 80, false},
		{":::443", "::", 443, false},
		{"127.0.0.1:8080", "127.0.0.1", 8080, false},
		{"invalid", "", 0, true},
		{"0.0.0.0:99999", "", 0, true},
	}
	for _, tt := range tests {
		host, port, err := ParseAddress(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			if host != tt.host || port != tt.port {
				t.Errorf("ParseAddress(%q) = (%q, %d), want (%q, %d)", tt.input, host, port, tt.host, tt.port)
			}
		}
	}
}

func TestDiff(t *testing.T) {
	prev := []Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 443},
	}
	curr := []Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 443},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 8080},
	}

	added, removed := Diff(prev, curr)

	if len(added) != 1 || added[0].Port != 8080 {
		t.Errorf("expected added port 8080, got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 80 {
		t.Errorf("expected removed port 80, got %v", removed)
	}
}

func TestDiffNoChanges(t *testing.T) {
	ports := []Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 22},
	}
	added, removed := Diff(ports, ports)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
