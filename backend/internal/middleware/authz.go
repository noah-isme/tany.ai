package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

const adminRole = "admin"

// AuthzAdmin ensures the authenticated principal carries the admin role.
func AuthzAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "authentication required", nil)
			return
		}

		if !hasRole(claims.Roles, adminRole) {
			httpapi.RespondError(c, http.StatusForbidden, httpapi.ErrorCodeForbidden, "admin access required", nil)
			return
		}

		c.Next()
	}
}

func hasRole(roles []string, desired string) bool {
	for _, role := range roles {
		if strings.EqualFold(role, desired) {
			return true
		}
	}
	return false
}
