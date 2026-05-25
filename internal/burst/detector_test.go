package burst

import (
	"testing"
	"time"
)

func TestNewInvalidWindow(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewInvalidThreshold(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewValid(t *testing.T) {
	d, err := New(time.Second, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestNoBurstBelowThreshold(t *testing.T) {
	d, _ := New(time.Second, 5)
	now := time.Now()
	for i := 0; i < 5; i++ {
		if d.Record(now) {
			t.Fatalf("false positive burst at i=%d", i)
		}
	}
}

func TestBurstDetectedAboveThreshold(t *testing.T) {
	d, _ := New(time.Second, 3)
	now := time.Now()
	d.Record(now)
	d.Record(now)
	d.Record(now)
	if !d.Record(now) {
		t.Fatal("expected burst to be detected on 4th record")
	}
}

func TestWindowEvictsOldEntries(t *testing.T) {
	d, _ := New(500*time.Millisecond, 2)
	old := time.Now().Add(-time.Second)
	d.Record(old)
	d.Record(old)
	// old entries should be outside the window; no burst expected
	now := time.Now()
	if d.Record(now) {
		t.Fatal("expected old entries to be evicted")
	}
}

func TestCountReflectsWindow(t *testing.T) {
	d, _ := New(time.Second, 10)
	now := time.Now()
	d.Record(now)
	d.Record(now)
	if got := d.Count(now); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestResetClearsState(t *testing.T) {
	d, _ := New(time.Second, 2)
	now := time.Now()
	d.Record(now)
	d.Record(now)
	d.Reset()
	if got := d.Count(now); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
}
