package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONLogger logs request metadata in structured JSON format.
func JSONLogger(route string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		entry := map[string]any{
			"route":         route,
			"status":        c.Writer.Status(),
			"method":        c.Request.Method,
			"path":          path,
			"ip":            c.ClientIP(),
			"latency_ms":    duration.Milliseconds(),
			"user_agent":    c.Request.UserAgent(),
			"cache_hit":     valueOrNil(c, "kb_cache_hit"),
			"prompt_length": valueOrNil(c, "prompt_length"),
			"chat_id":       valueOrNil(c, "chat_id"),
			"model":         valueOrNil(c, "model"),
		}

		payload, err := json.Marshal(entry)
		if err != nil {
			log.Printf("{\"route\":%q,\"error\":%q}", route, err.Error())
			return
		}
		log.Println(string(payload))
	}
}

func valueOrNil(c *gin.Context, key string) any {
	v, exists := c.Get(key)
	if !exists {
		return nil
	}
	return v
}
