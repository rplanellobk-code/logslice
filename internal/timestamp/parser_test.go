package timestamp

import (
	"testing"
	"time"
)

func TestParserAutoDetect(t *testing.T) {
	p := NewParser("", time.UTC)

	cases := []struct {
		line string
		wantYear int
		wantMonth time.Month
		wantDay int
	}{
		{"2024-03-15T08:30:00Z INFO server started", 2024, time.March, 15},
		{"2024-03-15T08:30:00.123456Z WARN high memory", 2024, time.March, 15},
		{"2024-03-15 08:30:00 ERROR disk full", 2024, time.March, 15},
		{"2024-03-15 08:30:00.000 DEBUG trace", 2024, time.March, 15},
		{"15/Mar/2024:08:30:00 +0000 GET /api", 2024, time.March, 15},
	}

	for _, tc := range cases {
		t.Run(tc.line[:20], func(t *testing.T) {
			got, err := p.Parse(tc.line)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tc.line, err)
			}
			if got.Year() != tc.wantYear || got.Month() != tc.wantMonth || got.Day() != tc.wantDay {
				t.Errorf("Parse(%q) = %v, want %d-%02d-%02d",
					tc.line, got, tc.wantYear, tc.wantMonth, tc.wantDay)
			}
		})
	}
}

func TestParserExplicitFormat(t *testing.T) {
	p := NewParser("2006-01-02 15:04:05", time.UTC)

	line := "2024-06-01 12:00:00 INFO application boot"
	got, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Hour() != 12 || got.Minute() != 0 {
		t.Errorf("unexpected time: %v", got)
	}
}

func TestParserUnrecognized(t *testing.T) {
	p := NewParser("", time.UTC)
	_, err := p.Parse("this line has no timestamp at all")
	if err == nil {
		t.Error("expected error for unrecognized timestamp, got nil")
	}
}

func TestParserNilLocation(t *testing.T) {
	p := NewParser("", nil)
	if p.loc != time.UTC {
		t.Error("expected UTC as default location")
	}
}
