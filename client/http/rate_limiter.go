package http

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 800 // Default to 800 requests per minute
	}

	// Set up a limiter that allows bursts up to the full minute limit
	// The rate is set to requestsPerMinute/60 per second, but the burst capacity
	// is set to the full minute's worth of requests. This means:
	// - All requests up to the minute limit will be processed immediately
	// - After the burst capacity is used, tokens refill at a rate of requestsPerMinute/60 per second
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), requestsPerMinute),
	}
}

func (r *RateLimiter) Take() {
	ctx := context.Background()
	r.limiter.Wait(ctx)
}
