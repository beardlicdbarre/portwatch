package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// Interval is how often to scan for port changes.
	Interval time.Duration `json:"interval"`
	// AlertFile is the path to write alerts to. Empty means stdout.
	AlertFile string `json:"alert_file"`
	// IgnorePorts is a list of ports to exclude from monitoring.
	IgnorePorts []int `json:"ignore_ports"`
	// Protocols is the list of protocols to monitor ("tcp", "udp").
	Protocols []string `json:"protocols"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval:    30 * time.Second,
		AlertFile:   "",
		IgnorePorts: []int{},
		Protocols:   []string{"tcp", "udp"},
	}
}

// Load reads a JSON config file from the given path.
// If path is empty, DefaultConfig is returned.
func Load(path string) (*Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	dec := json.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the config values are sensible.
func (c *Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("config: interval must be positive, got %v", c.Interval)
	}
	for _, p := range c.Protocols {
		if p != "tcp" && p != "udp" {
			return fmt.Errorf("config: unsupported protocol %q (must be tcp or udp)", p)
		}
	}
	return nil
}

// IsPortIgnored returns true if the given port number is in the ignore list.
func (c *Config) IsPortIgnored(port int) bool {
	for _, p := range c.IgnorePorts {
		if p == port {
			return true
		}
	}
	return false
}
