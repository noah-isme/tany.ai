package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLogger logs request metadata for specific routes in structured JSON format.
func JSONLogger(route string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		slog.Info("request",
			"route", route,
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"latency_ms", time.Since(start).Milliseconds(),
			"user_agent", c.Request.UserAgent(),
			"cache_hit", valueOrNil(c, "kb_cache_hit"),
			"prompt_length", valueOrNil(c, "prompt_length"),
			"chat_id", valueOrNil(c, "chat_id"),
			"model", valueOrNil(c, "model"),
		)
	}
}

func valueOrNil(c *gin.Context, key string) any {
	v, exists := c.Get(key)
	if !exists {
		return nil
	}
	return v
}
