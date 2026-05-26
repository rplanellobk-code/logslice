package linenum_test

import (
	"testing"
	"time"

	"github.com/aurc/logslice/internal/linenum"
	"github.com/aurc/logslice/internal/reader"
)

func makeLine(raw string) *reader.LogLine {
	return &reader.LogLine{Raw: raw, Timestamp: time.Now()}
}

func TestAnnotateNilLine(t *testing.T) {
	a := linenum.New()
	if err := a.Annotate(nil); err != linenum.ErrNilLine {
		t.Fatalf("expected ErrNilLine, got %v", err)
	}
}

func TestAnnotateIncrementsCounter(t *testing.T) {
	a := linenum.New()
	for i := 1; i <= 5; i++ {
		_ = a.Annotate(makeLine("hello"))
		if got := a.Count(); got != int64(i) {
			t.Fatalf("after %d calls Count() = %d, want %d", i, got, i)
		}
	}
}

func TestAnnotateDefaultPrefix(t *testing.T) {
	a := linenum.New()
	line := makeLine("world")
	_ = a.Annotate(line)
	want := "1 | world"
	if line.Raw != want {
		t.Fatalf("Raw = %q, want %q", line.Raw, want)
	}
}

func TestAnnotateCustomFormat(t *testing.T) {
	a := linenum.New(linenum.WithFormat("[%04d] "))
	line := makeLine("msg")
	_ = a.Annotate(line)
	want := "[0001] msg"
	if line.Raw != want {
		t.Fatalf("Raw = %q, want %q", line.Raw, want)
	}
}

func TestAnnotateEmptyFormatFallsBack(t *testing.T) {
	// Passing an empty string should keep the default format.
	a := linenum.New(linenum.WithFormat(""))
	line := makeLine("x")
	_ = a.Annotate(line)
	want := "1 | x"
	if line.Raw != want {
		t.Fatalf("Raw = %q, want %q", line.Raw, want)
	}
}

func TestReset(t *testing.T) {
	a := linenum.New()
	_ = a.Annotate(makeLine("a"))
	_ = a.Annotate(makeLine("b"))
	a.Reset()
	if got := a.Count(); got != 0 {
		t.Fatalf("Count after Reset = %d, want 0", got)
	}
	line := makeLine("c")
	_ = a.Annotate(line)
	want := "1 | c"
	if line.Raw != want {
		t.Fatalf("Raw after Reset = %q, want %q", line.Raw, want)
	}
}

func TestAnnotatePreservesTimestamp(t *testing.T) {
	a := linenum.New()
	now := time.Now()
	line := &reader.LogLine{Raw: "ts-check", Timestamp: now}
	_ = a.Annotate(line)
	if !line.Timestamp.Equal(now) {
		t.Fatal("Annotate must not modify the Timestamp field")
	}
}
