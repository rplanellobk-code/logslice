package burst

import (
	"testing"
	"time"
)

// BenchmarkRecordNoEviction measures Record when all entries stay within
// the window and no eviction occurs.
func BenchmarkRecordNoEviction(b *testing.B) {
	d, _ := New(time.Hour, b.N+1)
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Record(now)
	}
}

// BenchmarkRecordWithEviction measures Record when the window is very narrow
// so every call evicts the previous entry.
func BenchmarkRecordWithEviction(b *testing.B) {
	d, _ := New(time.Nanosecond, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Record(time.Now())
	}
}

// BenchmarkCountOnly measures Count without recording.
func BenchmarkCountOnly(b *testing.B) {
	d, _ := New(time.Second, 100)
	now := time.Now()
	for i := 0; i < 50; i++ {
		d.Record(now)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Count(now)
	}
}
