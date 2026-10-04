package ratelimiter

import (
	"context"
	"sync"
	"time"
)

// KeyedLimiter gives every key its own sliding window, and caps all keys together with one more.
// A request goes out only when both have room, and is recorded in both with the same timestamp:
// taking the key's slot first and then blocking on the total would stamp the key before the
// request actually leaves, and let the key exceed its budget as the server counts it.
type KeyedLimiter struct {
	window    time.Duration
	limitFor  func(key string) int
	total     slidingWindow
	keys      map[string]*slidingWindow
	lastSweep time.Time
	mu        sync.Mutex
}

// NewKeyedLimiter creates a limiter that allows limitFor(key) requests per window for each key,
// and totalRequests per window across all keys.
func NewKeyedLimiter(totalRequests int, window time.Duration, limitFor func(key string) int) *KeyedLimiter {
	return &KeyedLimiter{
		window:   window,
		limitFor: limitFor,
		total: slidingWindow{
			maxRequests: totalRequests,
			window:      window,
		},
		keys: map[string]*slidingWindow{},
	}
}

// Allow returns true if a request with this key can be made immediately.
// If true, the request is recorded and counts toward both the key's and the total limit.
func (l *KeyedLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.sweep(now)

	w := l.windowFor(key)
	w.cleanup(now)
	l.total.cleanup(now)

	if !w.hasRoom(now) || !l.total.hasRoom(now) {
		return false
	}

	w.record(now)
	l.total.record(now)

	return true
}

// Pause blocks requests with this key for at least d, extending an existing pause but never
// shortening it. Other keys are not affected.
func (l *KeyedLimiter) Pause(key string, d time.Duration) {
	if d <= 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.windowFor(key).pause(time.Now().Add(d))
}

// Wait blocks until a request with this key can be made, then records it.
// Returns an error if the context is canceled or times out.
func (l *KeyedLimiter) Wait(ctx context.Context, key string) error {
	return wait(ctx,
		func() bool { return l.Allow(key) },
		func() time.Duration { return l.nextAvailable(key) })
}

// nextAvailable returns how long to wait for the next request slot for this key
func (l *KeyedLimiter) nextAvailable(key string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	w := l.windowFor(key)
	w.cleanup(now)
	l.total.cleanup(now)

	return w.withPause(now, max(w.windowDelay(now), l.total.windowDelay(now)))
}

// Keys are per URL path, so a refresh creates one per resource. Dropping idle ones once per window
// keeps the map to roughly the keys used in the last window. A paused key is never idle: dropping
// it would silently lift the pause.
func (l *KeyedLimiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < l.window {
		return
	}

	l.lastSweep = now

	for key, w := range l.keys {
		w.cleanup(now)

		if w.idle(now) {
			delete(l.keys, key)
		}
	}
}

func (l *KeyedLimiter) windowFor(key string) *slidingWindow {
	w, ok := l.keys[key]
	if !ok {
		w = &slidingWindow{
			maxRequests: l.limitFor(key),
			window:      l.window,
		}
		l.keys[key] = w
	}

	return w
}
