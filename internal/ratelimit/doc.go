// Package ratelimit implements a token-bucket rate limiter designed for
// use in logslice processing pipelines.
//
// # Overview
//
// When processing very large log archives it can be desirable to limit
// throughput — for example to avoid saturating downstream storage or to
// replay logs at a controlled speed for debugging.
//
// # Usage
//
//	limiter := ratelimit.New(1000) // 1 000 lines/s
//	defer limiter.Stop()
//
//	for _, line := range lines {
//		if err := limiter.Wait(ctx); err != nil {
//			return err
//		}
//		process(line)
//	}
//
// Passing zero or a negative value to New disables rate limiting; Wait
// returns immediately in that case.
package ratelimit
