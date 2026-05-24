// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the throughput of log line processing in logslice pipelines.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which log lines are processed.
type Limiter struct {
	ticker  *time.Ticker
	tokens  chan struct{}
	cancel  context.CancelFunc
	done    chan struct{}
}

// New creates a Limiter that allows up to linesPerSecond lines per second.
// If linesPerSecond is zero or negative, no rate limiting is applied.
func New(linesPerSecond int) *Limiter {
	if linesPerSecond <= 0 {
		return &Limiter{}
	}

	interval := time.Second / time.Duration(linesPerSecond)
	ticker := time.NewTicker(interval)
	tokens := make(chan struct{}, linesPerSecond)
	ctx, cancel := context.WithCancel(context.Background())

	l := &Limiter{
		ticker: ticker,
		tokens: tokens,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	go l.produce(ctx)
	return l
}

// produce fills the token channel on each tick.
func (l *Limiter) produce(ctx context.Context) {
	defer close(l.done)
	for {
		select {
		case <-ctx.Done():
			l.ticker.Stop()
			return
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket full; drop token
			}
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
// Returns ctx.Err() if the context expires before a token is acquired.
// If the Limiter was created with no rate limit, Wait returns immediately.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.tokens == nil {
		return nil
	}
	select {
		case <-l.tokens:
			return nil
		case <-ctx.Done():
			return ctx.Err()
	}
}

// Stop shuts down the internal ticker goroutine.
// It is safe to call Stop on an unlimited Limiter.
func (l *Limiter) Stop() {
	if l.cancel == nil {
		return
	}
	l.cancel()
	<-l.done
}
