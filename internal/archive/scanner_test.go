package archive_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/archive"
)

func TestNewScannerEmptyPattern(t *testing.T) {
	_, err := archive.NewScanner("")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestScanDir(t *testing.T) {
	dir := t.TempDir()
	names := []string{"b.log", "a.log", "c.log"}
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("data"), 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}
	// add a sub-directory that should be skipped
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatalf("setup mkdir: %v", err)
	}

	s, err := archive.NewScanner(dir)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}
	files, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}
	// must be sorted
	if filepath.Base(files[0].Path) != "a.log" {
		t.Errorf("expected a.log first, got %s", files[0].Path)
	}
}

func TestScanGlob(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"app-1.log", "app-2.log", "other.txt"} {
		os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o644)
	}

	pattern := filepath.Join(dir, "app-*.log")
	s, err := archive.NewScanner(pattern)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}
	files, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestScanGlobNoMatch(t *testing.T) {
	dir := t.TempDir()
	s, _ := archive.NewScanner(filepath.Join(dir, "*.log"))
	files, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestFileSizePopulated(t *testing.T) {
	dir := t.TempDir()
	content := []byte("hello world")
	os.WriteFile(filepath.Join(dir, "test.log"), content, 0o644)

	s, _ := archive.NewScanner(dir)
	files, _ := s.Scan()
	if files[0].Size != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), files[0].Size)
	}
}
