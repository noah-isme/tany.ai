package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS adds Cross-Origin Resource Sharing headers to responses.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Build allowed origins from environment
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}
		
		// Add production frontend origin from environment
		if frontendOrigin := os.Getenv("FRONTEND_ORIGIN"); frontendOrigin != "" {
			// Support multiple origins separated by comma
			origins := strings.Split(frontendOrigin, ",")
			for _, o := range origins {
				o = strings.TrimSpace(o)
				if o != "" {
					allowedOrigins = append(allowedOrigins, o)
				}
			}
		}

		// For API clients and development tools that don't send Origin
		if origin == "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "43200") // 12 hours
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
			return
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || strings.HasPrefix(origin, allowedOrigin) {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Expose-Headers", "Content-Length")
			c.Header("Access-Control-Max-Age", "43200") // 12 hours

			// Handle preflight requests
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}