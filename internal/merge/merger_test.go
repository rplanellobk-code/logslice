package merge

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/reader"
)

func makeChan(lines []reader.LogLine) <-chan reader.LogLine {
	ch := make(chan reader.LogLine, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func ts(sec int) time.Time {
	return time.Unix(int64(sec), 0).UTC()
}

func TestNewNoSources(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil sources")
	}
	_, err = New([]<-chan reader.LogLine{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestMergeSingleSource(t *testing.T) {
	lines := []reader.LogLine{
		{Raw: "a", Timestamp: ts(1)},
		{Raw: "b", Timestamp: ts(2)},
		{Raw: "c", Timestamp: ts(3)},
	}
	m, err := New([]<-chan reader.LogLine{makeChan(lines)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got []reader.LogLine
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
}

func TestMergeMultipleSources(t *testing.T) {
	src1 := []reader.LogLine{
		{Raw: "s1-1", Timestamp: ts(1)},
		{Raw: "s1-3", Timestamp: ts(3)},
		{Raw: "s1-5", Timestamp: ts(5)},
	}
	src2 := []reader.LogLine{
		{Raw: "s2-2", Timestamp: ts(2)},
		{Raw: "s2-4", Timestamp: ts(4)},
		{Raw: "s2-6", Timestamp: ts(6)},
	}
	m, _ := New([]<-chan reader.LogLine{makeChan(src1), makeChan(src2)})
	var got []reader.LogLine
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 6 {
		t.Fatalf("want 6 lines, got %d", len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i].Timestamp.Before(got[i-1].Timestamp) {
			t.Errorf("out of order at index %d: %v before %v", i, got[i].Timestamp, got[i-1].Timestamp)
		}
	}
}

func TestMergeEmptySource(t *testing.T) {
	empty := make(chan reader.LogLine)
	close(empty)
	lines := []reader.LogLine{{Raw: "only", Timestamp: ts(1)}}
	m, _ := New([]<-chan reader.LogLine{empty, makeChan(lines)})
	var got []reader.LogLine
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 line, got %d", len(got))
	}
}
