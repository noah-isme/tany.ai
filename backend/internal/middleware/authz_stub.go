package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

const (
	envEnableAdminGuard = "ENABLE_ADMIN_GUARD"
	envAdminGuardMode   = "ADMIN_GUARD_MODE"
	envAppEnv           = "APP_ENV"
)

// AuthzAdminStub provides a toggleable middleware placeholder for future auth implementation.
func AuthzAdminStub() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !guardEnabled() {
			c.Next()
			return
		}

		mode := strings.ToLower(os.Getenv(envAdminGuardMode))
		switch mode {
		case "forbidden", "403":
			httpapi.RespondError(c, http.StatusForbidden, httpapi.ErrorCodeForbidden, "admin access is disabled", nil)
			return
		default:
			httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "admin authentication required", nil)
			return
		}
	}
}

func guardEnabled() bool {
	if value := strings.ToLower(os.Getenv(envEnableAdminGuard)); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return strings.ToLower(os.Getenv(envAppEnv)) == "prod"
}
