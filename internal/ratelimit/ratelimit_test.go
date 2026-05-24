package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/ratelimit"
)

// TestNewUnlimitedWaitReturnsImmediately verifies that a zero-rate limiter
// never blocks.
func TestNewUnlimitedWaitReturnsImmediately(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 1000; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("unlimited limiter blocked for %v, expected near-zero", elapsed)
	}
}

// TestRateLimiterThrottles verifies that a low-rate limiter slows processing.
func TestRateLimiterThrottles(t *testing.T) {
	const rate = 200 // lines per second
	l := ratelimit.New(rate)
	defer l.Stop()

	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	elapsed := time.Since(start)
	// 5 tokens at 200/s should take roughly 25 ms; allow generous window.
	if elapsed < 10*time.Millisecond {
		t.Errorf("limiter did not throttle: elapsed %v", elapsed)
	}
}

// TestWaitCancelledContext verifies that Wait honours context cancellation.
func TestWaitCancelledContext(t *testing.T) {
	// Very low rate so the token bucket is almost never filled.
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the first token immediately so the bucket is empty.
	_ = l.Wait(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

// TestStopIsIdempotent ensures Stop can be called multiple times safely.
func TestStopIsIdempotent(t *testing.T) {
	l := ratelimit.New(100)
	l.Stop()
	l.Stop() // must not panic
}

// TestNegativeRateIsUnlimited mirrors the zero-rate behaviour.
func TestNegativeRateIsUnlimited(t *testing.T) {
	l := ratelimit.New(-5)
	defer l.Stop()

	if err := l.Wait(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
