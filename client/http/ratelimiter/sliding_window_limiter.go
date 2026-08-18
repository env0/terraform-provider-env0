package ratelimiter

import (
	"context"
	"math/rand/v2"
	"sync"
	"time"
)

// Spread applied to the wake-up of requests waiting out a pause. A pause is one shared deadline,
// so without it every waiter resumes in the same instant and hits the server as one burst - the
// behaviour the caller's jittered backoff is trying to avoid. Capped so a short pause stays short.
const (
	pauseWakeupSpread    = 0.25
	pauseWakeupSpreadMax = time.Second
)

// SlidingWindowLimiter implements sliding window rate limiting.
// It tracks exact request timestamps to enforce: maxRequests per window duration.
type SlidingWindowLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	// No request is allowed before this point in time, regardless of the window budget.
	// Set by Pause when the server pushes back (429).
	pausedUntil time.Time
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

	if now.Before(l.pausedUntil) {
		return false
	}

	if len(l.requests) < l.maxRequests {
		l.requests = append(l.requests, now)

		return true
	}

	return false
}

// Pause blocks all requests for at least d, extending an existing pause but never shortening it.
// A paused limiter keeps its window budget: requests resume as soon as the pause expires.
func (l *SlidingWindowLimiter) Pause(d time.Duration) {
	if d <= 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	until := time.Now().Add(d)
	if until.After(l.pausedUntil) {
		l.pausedUntil = until
	}
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

	delay := time.Duration(0)

	if len(l.requests) >= l.maxRequests {
		// Wait for the oldest request to expire
		oldest := l.requests[0]
		delay = oldest.Add(l.window).Sub(now)
	}

	if pause := l.pausedUntil.Sub(now); pause > delay {
		delay = pause + pauseWakeupJitter(pause)
	}

	return delay
}

// pauseWakeupJitter returns a random extra wait in [0, min(pause*pauseWakeupSpread, pauseWakeupSpreadMax)).
func pauseWakeupJitter(pause time.Duration) time.Duration {
	spread := min(time.Duration(float64(pause)*pauseWakeupSpread), pauseWakeupSpreadMax)
	if spread <= 0 {
		return 0
	}

	return rand.N(spread)
}
