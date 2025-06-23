package http

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket algorithm for rate limiting
type RateLimiter struct {
	tokens         int       // Current number of tokens
	maxTokens      int       // Maximum number of tokens (equal to requests per minute)
	lastRefillTime time.Time // Last time the bucket was refilled
	mutex          sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified requests per minute
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 800 // Default to 800 requests per minute
	}

	return &RateLimiter{
		tokens:         requestsPerMinute, // Start with a full bucket
		maxTokens:      requestsPerMinute,
		lastRefillTime: time.Now(),
		mutex:          sync.Mutex{},
	}
}

// Take attempts to take a token from the bucket, blocking if necessary
func (r *RateLimiter) Take() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if a minute has passed since last refill
	now := time.Now()
	elapsed := now.Sub(r.lastRefillTime)

	// If a minute or more has passed, fully refill the bucket
	if elapsed >= time.Minute {
		r.tokens = r.maxTokens
		r.lastRefillTime = now
	}

	// If no tokens available, wait until the minute is up
	if r.tokens <= 0 {
		// Calculate time remaining in the current minute window
		timeToNextMinute := time.Minute - now.Sub(r.lastRefillTime)
		if timeToNextMinute < 0 {
			timeToNextMinute = 0
		}

		// Release lock while waiting
		r.mutex.Unlock()
		time.Sleep(timeToNextMinute)
		r.mutex.Lock()

		// Refill bucket after waiting
		r.tokens = r.maxTokens
		r.lastRefillTime = time.Now()
	}

	// Take a token
	r.tokens--
}
