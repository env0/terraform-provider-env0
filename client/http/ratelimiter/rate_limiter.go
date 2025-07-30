package ratelimiter

import "context"

type RateLimiter interface {
	Allow() bool
	Wait(ctx context.Context) error
}
