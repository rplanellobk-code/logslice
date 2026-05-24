package sampler

import (
	"testing"
)

func TestNewInvalidN(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewValidN(t *testing.T) {
	s, err := New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %d", s.Rate())
	}
}

func TestKeepEveryLine(t *testing.T) {
	s, _ := New(1)
	for i := 0; i < 10; i++ {
		if !s.Keep(Line{}) {
			t.Fatalf("line %d should be kept with n=1", i)
		}
	}
}

func TestKeepEveryNthLine(t *testing.T) {
	s, _ := New(3)
	kept := 0
	for i := 0; i < 9; i++ {
		if s.Keep(Line{}) {
			kept++
		}
	}
	if kept != 3 {
		t.Fatalf("expected 3 kept lines, got %d", kept)
	}
}

func TestKeepFirstIsAlwaysKept(t *testing.T) {
	s, _ := New(5)
	if !s.Keep(Line{}) {
		t.Fatal("first line should always be kept")
	}
}

func TestReset(t *testing.T) {
	s, _ := New(2)
	s.Keep(Line{}) // counter=1 kept
	s.Keep(Line{}) // counter=2 skipped
	s.Reset()
	if !s.Keep(Line{}) {
		t.Fatal("first line after reset should be kept")
	}
}

func TestConcurrentKeep(t *testing.T) {
	s, _ := New(2)
	done := make(chan struct{})
	for i := 0; i < 4; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				s.Keep(Line{})
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 4; i++ {
		<-done
	}
}
