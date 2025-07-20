package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RateLimiterMiddleware returns an Echo middleware that uses RateLimiter
func RateLimiterMiddleware(rl *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rl.Allow() {
				return next(c)
			}
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
			})
		}
	}
}
