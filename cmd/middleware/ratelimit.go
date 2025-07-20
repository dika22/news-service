package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// ===== Middleware Per IP =====
type IPBasedLimiter struct {
	limiters map[string]*RateLimiter
	lock         sync.Mutex
	maxTokens     int
	interval      time.Duration
	jitterPercent float64
}

func RateLimiterMiddleware(maxTokens int, interval time.Duration, jitter float64) *IPBasedLimiter {
	return &IPBasedLimiter{
		limiters:      make(map[string]*RateLimiter),
		maxTokens:     maxTokens,
		interval:      interval,
		jitterPercent: jitter,
	}
}

func (m *IPBasedLimiter) getLimiter(ip string) *RateLimiter {
	m.lock.Lock()
	defer m.lock.Unlock()

	if limiter, ok := m.limiters[ip]; ok {
		return limiter
	}

	// New limiter for new IP
	limiter := NewRateLimiter(m.maxTokens, m.interval, m.jitterPercent)
	m.limiters[ip] = limiter
	return limiter
}

func (m *IPBasedLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := m.getLimiter(ip)

			if limiter.Allow() {
				return next(c)
			}
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
				"ip":    ip,
			})
		}
	}
}
