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
	go rl.fill()
	return rl
}

func (rl *RateLimiter) fill() {
	for {
		select {
		case <-rl.stop:
			return
		default:
			delay := rl.jitteredInterval()
			time.Sleep(delay)

			select {
			case rl.tokens <- struct{}{}:
			default:
				// buffer full
			}
		}
	}
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

func (rl *RateLimiter) jitteredInterval() time.Duration {
	jitter := 1 + (rand.Float64()*rl.jitterPercent*2 - rl.jitterPercent)
	return time.Duration(float64(rl.interval) * jitter)
}
