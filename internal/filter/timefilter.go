package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/reader"
)

// TimeFilter holds the start and end bounds for log line filtering.
type TimeFilter struct {
	Start time.Time
	End   time.Time
}

// NewTimeFilter creates a TimeFilter from start and end times.
// Either boundary may be zero to indicate an open-ended range.
func NewTimeFilter(start, end time.Time) (*TimeFilter, error) {
	if !start.IsZero() && !end.IsZero() && end.Before(start) {
		return nil, fmt.Errorf("filter: end time %v is before start time %v", end, start)
	}
	return &TimeFilter{Start: start, End: end}, nil
}

// Match reports whether the given LogLine falls within the filter's time range.
func (f *TimeFilter) Match(line reader.LogLine) bool {
	if !f.Start.IsZero() && line.Timestamp.Before(f.Start) {
		return false
	}
	if !f.End.IsZero() && line.Timestamp.After(f.End) {
		return false
	}
	return true
}

// Filter returns only the LogLines whose timestamps fall within the range.
func (f *TimeFilter) Filter(lines []reader.LogLine) []reader.LogLine {
	result := make([]reader.LogLine, 0, len(lines))
	for _, l := range lines {
		if f.Match(l) {
			result = append(result, l)
		}
	}
	return result
}
