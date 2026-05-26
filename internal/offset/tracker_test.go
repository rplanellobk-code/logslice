package offset_test

import (
	"io"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/offset"
)

func TestNewNilReader(t *testing.T) {
	_, err := offset.New(nil)
	if err == nil {
		t.Fatal("expected error for nil ReadSeeker")
	}
}

func TestInitialOffsetIsZero(t *testing.T) {
	tr, err := offset.New(strings.NewReader("hello\nworld\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := tr.Offset(); got != 0 {
		t.Fatalf("expected initial offset 0, got %d", got)
	}
}

func TestReadLineAdvancesOffset(t *testing.T) {
	input := "foo\nbar\n"
	tr, err := offset.New(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line, err := tr.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != "foo" {
		t.Fatalf("expected 'foo', got %q", line)
	}
	if got := tr.Offset(); got != 4 {
		t.Fatalf("expected offset 4 after first line, got %d", got)
	}
}

func TestReadLineEOF(t *testing.T) {
	tr, err := offset.New(strings.NewReader("only\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr.ReadLine() // consume "only"
	_, err = tr.ReadLine()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

func TestSeekToRestoresPosition(t *testing.T) {
	input := "line1\nline2\nline3\n"
	tr, err := offset.New(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read past line1
	tr.ReadLine()
	saved := tr.Offset() // should be 6
	tr.ReadLine()        // consume line2

	if err := tr.SeekTo(saved); err != nil {
		t.Fatalf("SeekTo failed: %v", err)
	}
	if got := tr.Offset(); got != saved {
		t.Fatalf("expected offset %d after seek, got %d", saved, got)
	}

	line, err := tr.ReadLine()
	if err != nil {
		t.Fatalf("unexpected error after seek: %v", err)
	}
	if line != "line2" {
		t.Fatalf("expected 'line2' after seek, got %q", line)
	}
}

func TestSeekToNegativeOffset(t *testing.T) {
	tr, err := offset.New(strings.NewReader("data\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tr.SeekTo(-1); err != offset.ErrNegativeOffset {
		t.Fatalf("expected ErrNegativeOffset, got %v", err)
	}
}
