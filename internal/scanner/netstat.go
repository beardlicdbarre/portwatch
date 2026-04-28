package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// NetstatScanner uses the system netstat command to list open ports.
type NetstatScanner struct{}

// NewNetstatScanner returns a new NetstatScanner.
func NewNetstatScanner() *NetstatScanner {
	return &NetstatScanner{}
}

// Scan executes netstat and parses the output into a slice of Ports.
func (s *NetstatScanner) Scan() ([]Port, error) {
	args := netstatArgs()
	cmd := exec.Command("netstat", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("netstat failed: %w", err)
	}
	return parseNetstatOutput(out)
}

func netstatArgs() []string {
	if runtime.GOOS == "darwin" {
		return []string{"-an", "-p", "tcp"}
	}
	// Linux
	return []string{"-tlnup"}
}

func parseNetstatOutput(data []byte) ([]Port, error) {
	var ports []Port
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		proto := strings.ToLower(fields[0])
		if !strings.HasPrefix(proto, "tcp") && !strings.HasPrefix(proto, "udp") {
			continue
		}
		// Local address is typically the 4th field (index 3)
		localAddr := fields[3]
		host, port, err := ParseAddress(localAddr)
		if err != nil {
			continue
		}
		// On Linux, netstat -tlnup includes a State column; skip non-listening rows.
		if runtime.GOOS != "darwin" && !isListeningState(fields) {
			continue
		}
		p := Port{
			Protocol: proto,
			Address:  host,
			Port:     port,
			State:    "LISTEN",
		}
		if !seen[p.Key()] {
			seen[p.Key()] = true
			ports = append(ports, p)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading netstat output: %w", err)
	}
	return ports, nil
}

// isListeningState checks whether a parsed netstat line represents a listening
// socket by looking for a "LISTEN" state field anywhere in the fields slice.
func isListeningState(fields []string) bool {
	for _, f := range fields {
		if strings.EqualFold(f, "listen") {
			return true
		}
	}
	return false
}
