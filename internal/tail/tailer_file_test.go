package tail_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/tail"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestTailFileIntegration(t *testing.T) {
	var sb strings.Builder
	for i := 1; i <= 50; i++ {
		sb.WriteString("log line\n")
	}
	sb.WriteString("second to last\n")
	sb.WriteString("last line\n")

	path := writeTempFile(t, sb.String())

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	tl, err := tail.New(2)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines, err := tl.Read(f)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "second to last" {
		t.Errorf("lines[0] = %q, want %q", lines[0], "second to last")
	}
	if lines[1] != "last line" {
		t.Errorf("lines[1] = %q, want %q", lines[1], "last line")
	}
}
