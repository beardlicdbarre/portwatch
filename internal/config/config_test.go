package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if len(cfg.Protocols) != 2 {
		t.Errorf("expected 2 default protocols, got %d", len(cfg.Protocols))
	}
	if cfg.AlertFile != "" {
		t.Errorf("expected empty alert file, got %q", cfg.AlertFile)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadValidFile(t *testing.T) {
	data := map[string]interface{}{
		"interval":     "10s",
		"alert_file":   "/tmp/alerts.log",
		"ignore_ports": []int{22, 80},
		"protocols":    []string{"tcp"},
	}
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.Interval)
	}
	if cfg.AlertFile != "/tmp/alerts.log" {
		t.Errorf("unexpected alert file: %q", cfg.AlertFile)
	}
	if len(cfg.IgnorePorts) != 2 {
		t.Errorf("expected 2 ignored ports, got %d", len(cfg.IgnorePorts))
	}
}

func TestValidateInvalidInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative interval")
	}
}

func TestValidateInvalidProtocol(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Protocols = []string{"tcp", "icmp"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unsupported protocol")
	}
}

func TestIsPortIgnored(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IgnorePorts = []int{22, 443, 8080}

	if !cfg.IsPortIgnored(22) {
		t.Error("expected port 22 to be ignored")
	}
	if cfg.IsPortIgnored(80) {
		t.Error("expected port 80 not to be ignored")
	}
}
