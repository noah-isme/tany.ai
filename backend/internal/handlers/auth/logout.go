package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
)

// Logout revokes the active refresh token and clears the cookie.
func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie(h.refreshCookieName)
	if err == nil && token != "" {
		hash := appauth.HashRefreshToken(token)
		_ = h.users.RevokeRefreshTokenByHash(c.Request.Context(), hash)
	}

	appauth.ClearRefreshCookie(c, h.refreshCookieName)
	c.Status(http.StatusNoContent)
}
