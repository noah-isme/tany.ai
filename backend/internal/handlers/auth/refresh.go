package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// Refresh rotates the refresh token cookie and returns a new access token.
func (h *Handler) Refresh(c *gin.Context) {
	rawToken, err := c.Cookie(h.refreshCookieName)
	if err != nil || rawToken == "" {
		unauthorizedRefresh(c)
		return
	}

	tokenHash := appauth.HashRefreshToken(rawToken)
	stored, err := h.users.FindRefreshTokenByHash(c.Request.Context(), tokenHash)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			appauth.ClearRefreshCookie(c, h.refreshCookieName)
			unauthorizedRefresh(c)
			return
		}
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}

	now := time.Now().UTC()
	if stored.ExpiresAt.Before(now) || stored.Revoked {
		_ = h.users.RevokeRefreshTokenByHash(c.Request.Context(), tokenHash)
		appauth.ClearRefreshCookie(c, h.refreshCookieName)
		unauthorizedRefresh(c)
		return
	}

	revoked, err := h.users.RevokeRefreshToken(c.Request.Context(), stored.ID)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}
	if !revoked {
		appauth.ClearRefreshCookie(c, h.refreshCookieName)
		unauthorizedRefresh(c)
		return
	}

	user, roles, err := h.users.GetByID(c.Request.Context(), stored.UserID)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			appauth.ClearRefreshCookie(c, h.refreshCookieName)
			unauthorizedRefresh(c)
			return
		}
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}

	accessToken, err := h.tokens.GenerateAccessToken(appauth.Subject{ID: user.ID, Email: user.Email, Roles: roles})
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}

	refreshToken, refreshHash, expiresAt, err := h.tokens.GenerateRefreshToken()
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}

	record := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: expiresAt.UTC(),
	}
	if err := h.users.CreateRefreshToken(c.Request.Context(), record); err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "refresh failed", nil)
		return
	}

	appauth.SetRefreshCookie(c, h.refreshCookieName, refreshToken, expiresAt)

	resp := authResponse{
		AccessToken: accessToken,
		User:        dto.NewUserResponse(user, roles),
	}
	c.JSON(http.StatusOK, resp)
}

func unauthorizedRefresh(c *gin.Context) {
	httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "invalid refresh token", nil)
}
