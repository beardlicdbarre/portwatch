package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto, addr string, port int) scanner.Port {
	return scanner.Port{
		Protocol: proto,
		Address:  addr,
		PortNum:  port,
	}
}

func TestNotifyAdded(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := scanner.DiffResult{
		Added: []scanner.Port{makePort("tcp", "0.0.0.0", 8080)},
	}

	if err := n.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "new port opened") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestNotifyRemoved(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	diff := scanner.DiffResult{
		Removed: []scanner.Port{makePort("tcp", "127.0.0.1", 22)},
	}

	if err := n.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected port 22 in output, got: %s", out)
	}
}

func TestNotifyNoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(&buf)

	if err := n.Notify(scanner.DiffResult{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNewNotifierDefaultsToStdout(t *testing.T) {
	n := alert.NewNotifier(nil)
	if n.Writer == nil {
		t.Error("expected non-nil writer when nil passed to NewNotifier")
	}
}
