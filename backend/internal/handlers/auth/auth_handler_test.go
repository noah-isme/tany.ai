package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLoginSuccess(t *testing.T) {
	repo := newUserRepoStub()
	passwordHash, err := appauth.HashPassword("Admin#12345")
	require.NoError(t, err)
	userID := uuid.New()
	repo.addUser(models.User{ID: userID, Email: "admin@example.com", PasswordHash: passwordHash}, []string{"admin"})

	tokenService, err := appauth.NewTokenService(strings.Repeat("a", 64), 15*time.Minute, 24*time.Hour)
	require.NoError(t, err)
	handler := NewHandler(repo, tokenService, appauth.NewRateLimiter(10, 10, time.Minute), "__Host_refresh")

	router := gin.New()
	router.POST("/login", handler.Login)

	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"Admin#12345"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.1:1234"
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp authResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.AccessToken)
	require.Equal(t, "admin@example.com", resp.User.Email)
	require.Contains(t, resp.User.Roles, "admin")
	require.NotEmpty(t, repo.lastCreatedRefresh.TokenHash)

	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	var refreshCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "__Host_refresh" {
			refreshCookie = c
			break
		}
	}
	require.NotNil(t, refreshCookie)
	require.True(t, refreshCookie.HttpOnly)
	require.True(t, refreshCookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, refreshCookie.SameSite)
	require.Equal(t, "/", refreshCookie.Path)
}

func TestLoginInvalidCredentials(t *testing.T) {
	repo := newUserRepoStub()
	passwordHash, err := appauth.HashPassword("Admin#12345")
	require.NoError(t, err)
	repo.addUser(models.User{ID: uuid.New(), Email: "admin@example.com", PasswordHash: passwordHash}, []string{"admin"})

	tokenService, err := appauth.NewTokenService(strings.Repeat("b", 64), 15*time.Minute, 24*time.Hour)
	require.NoError(t, err)
	handler := NewHandler(repo, tokenService, appauth.NewRateLimiter(10, 10, time.Minute), "__Host_refresh")

	router := gin.New()
	router.POST("/login", handler.Login)

	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "203.0.113.1:1234"
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Empty(t, repo.lastCreatedRefresh.TokenHash)
}

func TestLoginRateLimited(t *testing.T) {
	repo := newUserRepoStub()
	passwordHash, err := appauth.HashPassword("Admin#12345")
	require.NoError(t, err)
	repo.addUser(models.User{ID: uuid.New(), Email: "admin@example.com", PasswordHash: passwordHash}, []string{"admin"})

	tokenService, err := appauth.NewTokenService(strings.Repeat("c", 64), 15*time.Minute, 24*time.Hour)
	require.NoError(t, err)
	limiter := appauth.NewRateLimiter(1, 1, time.Minute)
	handler := NewHandler(repo, tokenService, limiter, "__Host_refresh")

	router := gin.New()
	router.POST("/login", handler.Login)

	doLogin := func() *httptest.ResponseRecorder {
		body := bytes.NewBufferString(`{"email":"admin@example.com","password":"Admin#12345"}`)
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "198.51.100.1:1234"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	require.Equal(t, http.StatusOK, doLogin().Code)
	require.Equal(t, http.StatusTooManyRequests, doLogin().Code)
}

func TestRefreshSuccess(t *testing.T) {
	repo := newUserRepoStub()
	passwordHash, err := appauth.HashPassword("Admin#12345")
	require.NoError(t, err)
	userID := uuid.New()
	repo.addUser(models.User{ID: userID, Email: "admin@example.com", PasswordHash: passwordHash}, []string{"admin"})

	tokenService, err := appauth.NewTokenService(strings.Repeat("d", 64), 15*time.Minute, 24*time.Hour)
	require.NoError(t, err)
	handler := NewHandler(repo, tokenService, appauth.NewRateLimiter(10, 10, time.Minute), "__Host_refresh")

	oldRawRefresh := "refresh-token"
	oldHash := appauth.HashRefreshToken(oldRawRefresh)
	repo.addRefreshToken(models.RefreshToken{ID: uuid.New(), UserID: userID, TokenHash: oldHash, ExpiresAt: time.Now().Add(time.Hour)})

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "__Host_refresh", Value: oldRawRefresh})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, repo.refreshTokens[oldHash].Revoked)
	require.NotEqual(t, oldHash, repo.lastCreatedRefresh.TokenHash)

	var resp authResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.AccessToken)
	require.Equal(t, "admin@example.com", resp.User.Email)
}

func TestRefreshInvalidToken(t *testing.T) {
	repo := newUserRepoStub()
	tokenService, err := appauth.NewTokenService(strings.Repeat("e", 64), 15*time.Minute, 24*time.Hour)
	require.NoError(t, err)
	handler := NewHandler(repo, tokenService, appauth.NewRateLimiter(10, 10, time.Minute), "__Host_refresh")

	router := gin.New()
	router.POST("/refresh", handler.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "__Host_refresh", Value: "unknown"})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLogoutClearsCookie(t *testing.T) {
	repo := newUserRepoStub()
	userID := uuid.New()
	repo.addUser(models.User{ID: userID, Email: "admin@example.com"}, []string{"admin"})
	raw := "refresh"
	hash := appauth.HashRefreshToken(raw)
	repo.addRefreshToken(models.RefreshToken{ID: uuid.New(), UserID: userID, TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour)})

	handler := NewHandler(repo, nil, nil, "__Host_refresh")
	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "__Host_refresh", Value: raw})
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.True(t, repo.refreshTokens[hash].Revoked)
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	found := false
	for _, c := range cookies {
		if c.Name == "__Host_refresh" {
			require.Equal(t, "", c.Value)
			require.Equal(t, -1, c.MaxAge)
			found = true
		}
	}
	require.True(t, found)
}

type userRepoStub struct {
	usersByEmail       map[string]userRecord
	usersByID          map[uuid.UUID]userRecord
	refreshTokens      map[string]models.RefreshToken
	refreshByID        map[uuid.UUID]string
	lastCreatedRefresh models.RefreshToken
}

type userRecord struct {
	user  models.User
	roles []string
}

func newUserRepoStub() *userRepoStub {
	return &userRepoStub{
		usersByEmail:  make(map[string]userRecord),
		usersByID:     make(map[uuid.UUID]userRecord),
		refreshTokens: make(map[string]models.RefreshToken),
		refreshByID:   make(map[uuid.UUID]string),
	}
}

func (s *userRepoStub) addUser(user models.User, roles []string) {
	record := userRecord{user: user, roles: append([]string(nil), roles...)}
	s.usersByEmail[strings.ToLower(user.Email)] = record
	s.usersByID[user.ID] = record
}

func (s *userRepoStub) addRefreshToken(token models.RefreshToken) {
	s.refreshTokens[token.TokenHash] = token
	s.refreshByID[token.ID] = token.TokenHash
}

func (s *userRepoStub) GetByEmail(ctx context.Context, email string) (models.User, []string, error) {
	rec, ok := s.usersByEmail[strings.ToLower(email)]
	if !ok {
		return models.User{}, nil, repos.ErrNotFound
	}
	return rec.user, append([]string(nil), rec.roles...), nil
}

func (s *userRepoStub) GetByID(ctx context.Context, id uuid.UUID) (models.User, []string, error) {
	rec, ok := s.usersByID[id]
	if !ok {
		return models.User{}, nil, repos.ErrNotFound
	}
	return rec.user, append([]string(nil), rec.roles...), nil
}

func (s *userRepoStub) CreateRefreshToken(ctx context.Context, token models.RefreshToken) error {
	s.refreshTokens[token.TokenHash] = token
	s.refreshByID[token.ID] = token.TokenHash
	s.lastCreatedRefresh = token
	return nil
}

func (s *userRepoStub) FindRefreshTokenByHash(ctx context.Context, hash string) (models.RefreshToken, error) {
	token, ok := s.refreshTokens[hash]
	if !ok {
		return models.RefreshToken{}, repos.ErrNotFound
	}
	return token, nil
}

func (s *userRepoStub) RevokeRefreshToken(ctx context.Context, id uuid.UUID) (bool, error) {
	hash, ok := s.refreshByID[id]
	if !ok {
		return false, nil
	}
	token := s.refreshTokens[hash]
	if token.Revoked {
		return false, nil
	}
	token.Revoked = true
	s.refreshTokens[hash] = token
	return true, nil
}

func (s *userRepoStub) RevokeRefreshTokenByHash(ctx context.Context, hash string) error {
	token, ok := s.refreshTokens[hash]
	if ok {
		token.Revoked = true
		s.refreshTokens[hash] = token
	}
	return nil
}

func (s *userRepoStub) DeleteExpiredRefreshTokens(ctx context.Context, before time.Time) (int64, error) {
	return 0, nil
}
