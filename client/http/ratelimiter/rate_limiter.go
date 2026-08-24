package ratelimiter

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
	// Pause holds back every request, not just the one that triggered it, for at least d.
	// Used when the server answers 429: the limit is server side, so sibling requests must
	// back off too instead of spending the rest of the window on responses that get blocked.
	Pause(d time.Duration)
}
