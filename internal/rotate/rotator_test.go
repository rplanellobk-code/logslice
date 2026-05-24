package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRotatorByLines(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir, "out", ByLines, 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	for i := 0; i < 7; i++ {
		if err := r.WriteLine(fmt.Sprintf("line %d", i)); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	r.Close()

	// 7 lines with threshold 3 → files 1,2,3 (3+3+1)
	if got := r.FilesCreated(); got != 3 {
		t.Errorf("FilesCreated = %d, want 3", got)
	}
	if got := r.LinesWritten(); got != 7 {
		t.Errorf("LinesWritten = %d, want 7", got)
	}
}

func TestRotatorBySize(t *testing.T) {
	dir := t.TempDir()
	// Each line "line X\n" is 7 bytes; threshold 20 → rotate after ~2 lines
	r, err := New(dir, "sz", BySize, 20)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	for i := 0; i < 6; i++ {
		if err := r.WriteLine(fmt.Sprintf("line %d", i)); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	r.Close()

	if r.FilesCreated() < 2 {
		t.Errorf("expected at least 2 files, got %d", r.FilesCreated())
	}
}

func TestRotatorFilesExist(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir, "log", ByLines, 2)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 4; i++ {
		_ = r.WriteLine("hello")
	}
	r.Close()

	for i := 1; i <= r.FilesCreated(); i++ {
		path := filepath.Join(dir, fmt.Sprintf("log_%04d.log", i))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}
}

func TestRotatorInvalidThreshold(t *testing.T) {
	dir := t.TempDir()
	_, err := New(dir, "x", ByLines, 0)
	if err == nil {
		t.Error("expected error for threshold=0, got nil")
	}
}

func TestRotatorLinesWrittenAccumulates(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir, "acc", ByLines, 10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	const total = 25
	for i := 0; i < total; i++ {
		if err := r.WriteLine("data"); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if got := r.LinesWritten(); got != total {
		t.Errorf("LinesWritten = %d, want %d", got, total)
	}
}
