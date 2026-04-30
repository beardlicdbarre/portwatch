package formatter

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func samplePort() scanner.Port {
	return scanner.Port{Proto: "tcp", Addr: "0.0.0.0", Port: 8080}
}

func TestFormatAddedText(t *testing.T) {
	f := New(FormatText, false)
	out := f.FormatAdded(samplePort())
	if !strings.Contains(out, "[ADDED]") {
		t.Errorf("expected [ADDED] in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestFormatRemovedText(t *testing.T) {
	f := New(FormatText, false)
	out := f.FormatRemoved(samplePort())
	if !strings.Contains(out, "[REMOVED]") {
		t.Errorf("expected [REMOVED] in output, got: %s", out)
	}
}

func TestFormatTextWithTimestamp(t *testing.T) {
	f := New(FormatText, true)
	out := f.FormatAdded(samplePort())
	// RFC3339 timestamps contain 'T' and 'Z'
	if !strings.Contains(out, "T") || !strings.Contains(out, "Z") {
		t.Errorf("expected RFC3339 timestamp in output, got: %s", out)
	}
}

func TestFormatJSON(t *testing.T) {
	f := New(FormatJSON, false)
	out := f.FormatAdded(samplePort())
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, `"event":"ADDED"`) {
		t.Errorf("expected event field in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"port":8080`) {
		t.Errorf("expected port field in JSON, got: %s", out)
	}
}

func TestFormatJSONWithTimestamp(t *testing.T) {
	f := New(FormatJSON, true)
	out := f.FormatRemoved(samplePort())
	if !strings.Contains(out, `"timestamp"`) {
		t.Errorf("expected timestamp field in JSON, got: %s", out)
	}
}

func TestFormatJSONNoTimestamp(t *testing.T) {
	f := New(FormatJSON, false)
	out := f.FormatAdded(samplePort())
	if strings.Contains(out, `"timestamp"`) {
		t.Errorf("unexpected timestamp field in JSON, got: %s", out)
	}
}
