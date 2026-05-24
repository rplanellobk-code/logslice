package linecount_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestNewInvalidBufSize(t *testing.T) {
	_, err := linecount.New(linecount.WithBufSize(0))
	if err == nil {
		t.Fatal("expected error for bufSize=0, got nil")
	}
}

func TestCountReaderEmpty(t *testing.T) {
	c, err := linecount.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	n, err := c.CountReader(strings.NewReader(""))
	if err != nil {
		t.Fatalf("CountReader: %v", err)
	}
	if n != 0 {
		t.Errorf("want 0, got %d", n)
	}
}

func TestCountReaderLines(t *testing.T) {
	c, _ := linecount.New()
	input := "line1\nline2\nline3\n"
	n, err := c.CountReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("CountReader: %v", err)
	}
	if n != 3 {
		t.Errorf("want 3, got %d", n)
	}
}

func TestCountReaderNoTrailingNewline(t *testing.T) {
	c, _ := linecount.New()
	n, err := c.CountReader(strings.NewReader("a\nb"))
	if err != nil {
		t.Fatalf("CountReader: %v", err)
	}
	if n != 2 {
		t.Errorf("want 2, got %d", n)
	}
}

func TestCountFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	content := "alpha\nbeta\ngamma\ndelta\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	c, _ := linecount.New()
	n, err := c.CountFile(path)
	if err != nil {
		t.Fatalf("CountFile: %v", err)
	}
	if n != 4 {
		t.Errorf("want 4, got %d", n)
	}
}

func TestCountFileMissing(t *testing.T) {
	c, _ := linecount.New()
	_, err := c.CountFile("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestCountReaderCustomBufSize(t *testing.T) {
	c, err := linecount.New(linecount.WithBufSize(128))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	n, err := c.CountReader(strings.NewReader("x\ny\nz\n"))
	if err != nil {
		t.Fatalf("CountReader: %v", err)
	}
	if n != 3 {
		t.Errorf("want 3, got %d", n)
	}
}
