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

// slidingWindow is the bookkeeping shared by SlidingWindowLimiter and KeyedLimiter. It holds no
// lock of its own: every method expects the caller to hold the lock that guards it.
type slidingWindow struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	// No request is allowed before this point in time, regardless of the window budget.
	// Set by Pause when the server pushes back (429).
	pausedUntil time.Time
}

// cleanup removes expired requests from the sliding window
func (w *slidingWindow) cleanup(now time.Time) {
	cutoff := now.Add(-w.window)

	// Find first request still within window
	i := 0
	for i < len(w.requests) && w.requests[i].Before(cutoff) {
		i++
	}

	// Remove expired requests
	if i > 0 {
		w.requests = w.requests[i:]
	}
}

// hasRoom reports whether a request can go out now. Call cleanup first.
func (w *slidingWindow) hasRoom(now time.Time) bool {
	return !now.Before(w.pausedUntil) && len(w.requests) < w.maxRequests
}

func (w *slidingWindow) record(now time.Time) {
	w.requests = append(w.requests, now)
}

// pause extends an existing pause but never shortens it.
func (w *slidingWindow) pause(until time.Time) {
	if until.After(w.pausedUntil) {
		w.pausedUntil = until
	}
}

// windowDelay returns how long until the window has a free slot, ignoring any pause. Call cleanup
// first.
func (w *slidingWindow) windowDelay(now time.Time) time.Duration {
	if len(w.requests) < w.maxRequests {
		return 0
	}

	// Wait for the oldest request to expire
	return w.requests[0].Add(w.window).Sub(now)
}

// idle reports whether dropping the window loses nothing: no request in it and no pause pending.
// Call cleanup first.
func (w *slidingWindow) idle(now time.Time) bool {
	return len(w.requests) == 0 && !now.Before(w.pausedUntil)
}

// withPause returns delay, or the remaining pause plus a wake-up jitter when the pause is longer.
func (w *slidingWindow) withPause(now time.Time, delay time.Duration) time.Duration {
	if pause := w.pausedUntil.Sub(now); pause > delay {
		return pause + pauseWakeupJitter(pause)
	}

	return delay
}

// SlidingWindowLimiter implements sliding window rate limiting.
// It tracks exact request timestamps to enforce: maxRequests per window duration.
type SlidingWindowLimiter struct {
	slidingWindow
	mu sync.Mutex
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter.
// maxRequests: maximum number of requests allowed
// window: time window duration
//
// Example: NewSlidingWindowLimiter(100, time.Hour) allows 100 requests per hour.
func NewSlidingWindowLimiter(maxRequests int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		slidingWindow: slidingWindow{
			maxRequests: maxRequests,
			window:      window,
		},
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

	if !l.hasRoom(now) {
		return false
	}

	l.record(now)

	return true
}

// Pause blocks all requests for at least d, extending an existing pause but never shortening it.
// A paused limiter keeps its window budget: requests resume as soon as the pause expires.
func (l *SlidingWindowLimiter) Pause(d time.Duration) {
	if d <= 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.pause(time.Now().Add(d))
}

// Wait blocks until a request can be made, then records it.
// Returns an error if the context is canceled or times out.
// Use this for blocking requests where you want to wait for rate limit clearance.
func (l *SlidingWindowLimiter) Wait(ctx context.Context) error {
	return wait(ctx, l.Allow, l.nextAvailable)
}

// nextAvailable returns how long to wait for the next request slot
func (l *SlidingWindowLimiter) nextAvailable() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanup(now)

	return l.withPause(now, l.windowDelay(now))
}

// wait retries allow until it succeeds, sleeping for nextAvailable in between.
// Returns an error if the context is canceled or times out.
func wait(ctx context.Context, allow func() bool, nextAvailable func() time.Duration) error {
	for {
		// Try to make the request immediately
		if allow() {
			return nil
		}

		// Calculate how long to wait
		delay := nextAvailable()
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

// pauseWakeupJitter returns a random extra wait in [0, min(pause*pauseWakeupSpread, pauseWakeupSpreadMax)).
func pauseWakeupJitter(pause time.Duration) time.Duration {
	spread := min(time.Duration(float64(pause)*pauseWakeupSpread), pauseWakeupSpreadMax)
	if spread <= 0 {
		return 0
	}

	return rand.N(spread)
}
