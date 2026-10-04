package ratelimiter

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(key string) bool
	Wait(ctx context.Context, key string) error
	// Pause holds back every request with this key, not just the one that triggered it, for at
	// least d. Used when the server answers 429: the WAF counts per method and path, so sibling
	// requests to the same key must back off too instead of spending the rest of the window on
	// responses that get blocked, while a 429 on one key says nothing about the others.
	Pause(key string, d time.Duration)
}
