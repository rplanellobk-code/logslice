package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNewStoreEmptyPath(t *testing.T) {
	_, err := checkpoint.NewStore("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	store, err := checkpoint.NewStore(tempPath(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	st, err := store.Load()
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if st.LastFile != "" || st.LastOffset != 0 {
		t.Errorf("expected zero State, got %+v", st)
	}
}

func TestSaveAndLoad(t *testing.T) {
	store, err := checkpoint.NewStore(tempPath(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	want := checkpoint.State{
		LastFile:      "/var/log/app/2024-01-15.log",
		LastOffset:    4096,
		LastTimestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.LastFile != want.LastFile {
		t.Errorf("LastFile: got %q, want %q", got.LastFile, want.LastFile)
	}
	if got.LastOffset != want.LastOffset {
		t.Errorf("LastOffset: got %d, want %d", got.LastOffset, want.LastOffset)
	}
	if !got.UpdatedAt.IsZero() == false && got.UpdatedAt.Before(time.Now()) == false {
		t.Errorf("UpdatedAt not set correctly: %v", got.UpdatedAt)
	}
}

func TestDelete(t *testing.T) {
	p := tempPath(t)
	store, _ := checkpoint.NewStore(p)
	_ = store.Save(checkpoint.State{LastFile: "x"})

	if err := store.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Delete")
	}
}

func TestDeleteMissingFileIsNoOp(t *testing.T) {
	store, _ := checkpoint.NewStore(tempPath(t))
	if err := store.Delete(); err != nil {
		t.Fatalf("Delete on missing file should be no-op, got: %v", err)
	}
}
