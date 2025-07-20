package middleware

import (
	"math/rand"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens        chan struct{}
	maxTokens     int
	interval      time.Duration
	jitterPercent float64
	stop          chan struct{}
	once          sync.Once
}

// NewRateLimiter initializes the rate limiter with pre-filled tokens.
func NewRateLimiter(maxTokens int, interval time.Duration, jitter float64) *RateLimiter {
	rl := &RateLimiter{
		tokens:        make(chan struct{}, maxTokens),
		maxTokens:     maxTokens,
		interval:      interval,
		jitterPercent: jitter,
		stop:          make(chan struct{}),
	}

	// Pre-fill tokens
	for i := 0; i < maxTokens; i++ {
		rl.tokens <- struct{}{}
	}

	go rl.fill()
	return rl
}

// fill refills tokens at interval ± jitter
func (rl *RateLimiter) fill() {
	ticker := time.NewTicker(rl.jitteredInterval())
	defer ticker.Stop()

	for {
		select {
		case <-rl.stop:
			return
		case <-ticker.C:
			// Refill 1 token if there's space
			select {
			case rl.tokens <- struct{}{}:
			default:
				// channel full, do nothing
			}
			ticker.Reset(rl.jitteredInterval())
		}
	}
}

// jitteredInterval returns interval ± jitter
func (rl *RateLimiter) jitteredInterval() time.Duration {
	jitter := 1 + (rand.Float64()*rl.jitterPercent*2 - rl.jitterPercent)
	return time.Duration(float64(rl.interval) * jitter)
}

// Allow checks if a token is available
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Stop stops the refill goroutine
func (rl *RateLimiter) Stop() {
	rl.once.Do(func() {
		close(rl.stop)
	})
}