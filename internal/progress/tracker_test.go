package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestTrackerSummary(t *testing.T) {
	tr := NewTracker(nil, 0)
	tr.AddBytes(2048)
	tr.AddLineIn()
	tr.AddLineIn()
	tr.AddLineOut()

	summary := tr.Summary()
	if !strings.Contains(summary, "read 2 lines") {
		t.Errorf("expected 'read 2 lines' in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "wrote 1 lines") {
		t.Errorf("expected 'wrote 1 lines' in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "2.0 KB") {
		t.Errorf("expected '2.0 KB' in summary, got: %s", summary)
	}
}

func TestTrackerStopPrintsFinal(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracker(&buf, 0)
	tr.AddLineIn()
	tr.AddLineOut()
	tr.Stop()

	out := buf.String()
	if !strings.Contains(out, "done:") {
		t.Errorf("expected 'done:' prefix in output, got: %s", out)
	}
}

func TestTrackerPeriodicReport(t *testing.T) {
	var buf bytes.Buffer
	tr := NewTracker(&buf, 20*time.Millisecond)
	tr.AddLineIn()
	time.Sleep(60 * time.Millisecond)
	tr.Stop()

	out := buf.String()
	if !strings.Contains(out, "progress:") {
		t.Errorf("expected at least one 'progress:' line, got: %s", out)
	}
}

func TestTrackerNilWriter(t *testing.T) {
	// Should not panic when writer is nil.
	tr := NewTracker(nil, 0)
	tr.AddBytes(512)
	tr.AddLineIn()
	tr.Stop()
}
