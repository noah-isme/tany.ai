package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

// RateLimitByIP applies a token bucket limiter keyed by client IP address.
func RateLimitByIP(limiter *appauth.RateLimiter) gin.HandlerFunc {
	if limiter == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		if key == "" {
			key = "unknown"
		}
		if !limiter.Allow(key) {
			httpapi.RespondError(c, http.StatusTooManyRequests, httpapi.ErrorCodeTooManyRequests, "too many requests", nil)
			return
		}
		c.Next()
	}
}
