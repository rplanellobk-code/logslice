package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/reader"
)

func makeLines(timestamps []string) []reader.LogLine {
	layout := "2006-01-02T15:04:05"
	lines := make([]reader.LogLine, 0, len(timestamps))
	for _, ts := range timestamps {
		t, _ := time.Parse(layout, ts)
		lines = append(lines, reader.LogLine{Timestamp: t, Raw: ts})
	}
	return lines
}

func TestNewTimeFilterInvalidRange(t *testing.T) {
	start := time.Now()
	end := start.Add(-time.Hour)
	_, err := filter.NewTimeFilter(start, end)
	if err == nil {
		t.Fatal("expected error for end before start, got nil")
	}
}

func TestTimeFilterBothBounds(t *testing.T) {
	start, _ := time.Parse("2006-01-02T15:04:05", "2024-01-01T10:00:00")
	end, _ := time.Parse("2006-01-02T15:04:05", "2024-01-01T12:00:00")
	f, err := filter.NewTimeFilter(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := makeLines([]string{
		"2024-01-01T09:59:59",
		"2024-01-01T10:00:00",
		"2024-01-01T11:00:00",
		"2024-01-01T12:00:00",
		"2024-01-01T12:00:01",
	})

	got := f.Filter(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestTimeFilterOpenStart(t *testing.T) {
	end, _ := time.Parse("2006-01-02T15:04:05", "2024-01-01T11:00:00")
	f, _ := filter.NewTimeFilter(time.Time{}, end)

	lines := makeLines([]string{
		"2024-01-01T09:00:00",
		"2024-01-01T11:00:00",
		"2024-01-01T12:00:00",
	})

	got := f.Filter(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestTimeFilterOpenEnd(t *testing.T) {
	start, _ := time.Parse("2006-01-02T15:04:05", "2024-01-01T11:00:00")
	f, _ := filter.NewTimeFilter(start, time.Time{})

	lines := makeLines([]string{
		"2024-01-01T10:00:00",
		"2024-01-01T11:00:00",
		"2024-01-01T13:00:00",
	})

	got := f.Filter(lines)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}
