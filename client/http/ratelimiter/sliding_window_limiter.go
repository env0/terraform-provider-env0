package ratelimiter

import (
	"context"
	"sync"
	"time"
)

// SlidingWindowLimiter implements sliding window rate limiting.
// It tracks exact request timestamps to enforce: maxRequests per window duration.
type SlidingWindowLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	mu          sync.Mutex
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter.
// maxRequests: maximum number of requests allowed
// window: time window duration
//
// Example: NewSlidingWindowLimiter(100, time.Hour) allows 100 requests per hour.
func NewSlidingWindowLimiter(maxRequests int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// Allow returns true if a request can be made immediately.
// If true, the request is recorded and counts toward the limit.
// Use this for non-blocking requests where you want to skip if rate limited.
func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanup(now)

	if len(l.requests) < l.maxRequests {
		l.requests = append(l.requests, now)

		return true
	}

	return false
}

// Wait blocks until a request can be made, then records it.
// Returns an error if the context is canceled or times out.
// Use this for blocking requests where you want to wait for rate limit clearance.
func (l *SlidingWindowLimiter) Wait(ctx context.Context) error {
	for {
		// Try to make the request immediately
		if l.Allow() {
			return nil
		}

		// Calculate how long to wait
		delay := l.nextAvailable()
		if delay <= 0 {
			continue // Should be available now, try again
		}

		// Wait for the delay or context cancellation
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			timer.Stop()
			// Try again after waiting
		case <-ctx.Done():
			timer.Stop()

			return ctx.Err()
		}
	}
}

// cleanup removes expired requests from the sliding window
func (l *SlidingWindowLimiter) cleanup(now time.Time) {
	cutoff := now.Add(-l.window)

	// Find first request still within window
	i := 0
	for i < len(l.requests) && l.requests[i].Before(cutoff) {
		i++
	}

	// Remove expired requests
	if i > 0 {
		l.requests = l.requests[i:]
	}
}

// nextAvailable returns how long to wait for the next request slot
func (l *SlidingWindowLimiter) nextAvailable() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanup(now)

	if len(l.requests) < l.maxRequests {
		return 0
	}

	// Wait for the oldest request to expire
	oldest := l.requests[0]

	return oldest.Add(l.window).Sub(now)
}
