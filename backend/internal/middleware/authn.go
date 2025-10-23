package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

const contextClaimsKey = "auth.claims"

// AccessTokenValidator defines the behaviour required to validate access tokens.
type AccessTokenValidator interface {
	ValidateAccessToken(token string) (*auth.Claims, error)
}

// Authn enforces JWT-based authentication on protected routes.
func Authn(validator AccessTokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "authentication required", nil)
			return
		}

		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "invalid authorization header", nil)
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "invalid authorization header", nil)
			return
		}

		claims, err := validator.ValidateAccessToken(token)
		if err != nil {
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "invalid or expired token", nil)
			return
		}

		c.Set(contextClaimsKey, claims)
		c.Next()
	}
}

// GetClaims retrieves access token claims from context.
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get(contextClaimsKey)
	if !ok {
		return nil, false
	}
	claims, ok := value.(*auth.Claims)
	return claims, ok
}
