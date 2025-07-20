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

func NewRateLimiter(maxTokens int, interval time.Duration, jitter float64) *RateLimiter {
	rl := &RateLimiter{
		tokens:        make(chan struct{}, maxTokens),
		maxTokens:     maxTokens,
		interval:      interval,
		jitterPercent: jitter,
		stop:          make(chan struct{}),
	}
	// Isi token awal
	for i := 0; i < maxTokens; i++ {
		rl.tokens <- struct{}{}
	}
	go rl.fill()
	return rl
}

func (rl *RateLimiter) fill() {
	ticker := time.NewTicker(rl.jitteredInterval())
	defer ticker.Stop()

	for {
		select {
		case <-rl.stop:
			return
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// full, skip
			}
			ticker.Reset(rl.jitteredInterval())
		}
	}
}

func (rl *RateLimiter) jitteredInterval() time.Duration {
	jitter := 1 + (rand.Float64()*rl.jitterPercent*2 - rl.jitterPercent)
	return time.Duration(float64(rl.interval) * jitter)
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Stop() {
	rl.once.Do(func() {
		close(rl.stop)
	})
}