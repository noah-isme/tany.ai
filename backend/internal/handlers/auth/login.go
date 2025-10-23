package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

var dummyPasswordHash = func() string {
	hash, err := appauth.HashPassword("dummy-padding-password")
	if err != nil {
		panic(err)
	}
	return hash
}()

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	AccessToken string           `json:"accessToken"`
	User        dto.UserResponse `json:"user"`
}

// Login authenticates an admin user using email/password credentials.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid payload", nil)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "email is required", nil)
		return
	}

	if h.limiter != nil {
		key := fmt.Sprintf("%s:%s", email, c.ClientIP())
		if !h.limiter.Allow(key) {
			httpapi.RespondError(c, http.StatusTooManyRequests, httpapi.ErrorCodeTooManyRequests, "too many login attempts", nil)
			return
		}
	}

	user, roles, err := h.users.GetByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, repos.ErrNotFound) {
			_ = appauth.ComparePassword(dummyPasswordHash, req.Password)
			respondInvalidCredentials(c)
			return
		}
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "login failed", nil)
		return
	}

	if err := appauth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		if errors.Is(err, appauth.ErrInvalidPassword) {
			respondInvalidCredentials(c)
			return
		}
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "login failed", nil)
		return
	}

	accessToken, err := h.tokens.GenerateAccessToken(appauth.Subject{ID: user.ID, Email: user.Email, Roles: roles})
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "login failed", nil)
		return
	}

	refreshToken, refreshHash, expiresAt, err := h.tokens.GenerateRefreshToken()
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "login failed", nil)
		return
	}

	record := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: expiresAt,
	}
	if err := h.users.CreateRefreshToken(c.Request.Context(), record); err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "login failed", nil)
		return
	}

	appauth.SetRefreshCookie(c, h.refreshCookieName, refreshToken, expiresAt)

	resp := authResponse{
		AccessToken: accessToken,
		User:        dto.NewUserResponse(user, roles),
	}
	c.JSON(http.StatusOK, resp)
}

func respondInvalidCredentials(c *gin.Context) {
	httpapi.RespondError(c, http.StatusUnauthorized, httpapi.ErrorCodeUnauthorized, "invalid credentials", nil)
}
