package auth

import (
	"time"

	"github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// TokenManager defines the token operations required by auth handlers.
type TokenManager interface {
	GenerateAccessToken(sub auth.Subject) (string, error)
	GenerateRefreshToken() (token string, hash string, expiresAt time.Time, err error)
}

// Handler groups authentication related HTTP handlers.
type Handler struct {
	users             repos.UserRepository
	tokens            TokenManager
	limiter           *auth.RateLimiter
	refreshCookieName string
}

// NewHandler constructs an auth Handler.
func NewHandler(users repos.UserRepository, tokens TokenManager, limiter *auth.RateLimiter, refreshCookieName string) *Handler {
	return &Handler{
		users:             users,
		tokens:            tokens,
		limiter:           limiter,
		refreshCookieName: refreshCookieName,
	}
}

// RefreshCookieName returns the configured refresh cookie name.
func (h *Handler) RefreshCookieName() string {
	return h.refreshCookieName
}
