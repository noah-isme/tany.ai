package auth_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/server"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAdminRoutesRequireAuthentication(t *testing.T) {
	srv, _, cleanup := newTestServer(t)
	t.Cleanup(cleanup)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", nil)
	rec := httptest.NewRecorder()

	srv.Engine().ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminRoutesRejectNonAdmin(t *testing.T) {
	srv, tokenService, cleanup := newTestServer(t)
	t.Cleanup(cleanup)

	token, err := tokenService.GenerateAccessToken(appauth.Subject{
		ID:    uuid.New(),
		Email: "user@example.com",
		Roles: []string{"user"},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	srv.Engine().ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminRoutesAllowAdmin(t *testing.T) {
	srv, tokenService, cleanup := newTestServer(t)
	t.Cleanup(cleanup)

	token, err := tokenService.GenerateAccessToken(appauth.Subject{
		ID:    uuid.New(),
		Email: "admin@example.com",
		Roles: []string{"admin"},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	srv.Engine().ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func newTestServer(t *testing.T) (*server.Server, *appauth.TokenService, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT value FROM embedding_config WHERE key = $1")).
		WithArgs("personalization").
		WillReturnError(sql.ErrNoRows)

	cfg := config.Config{
		AppEnv:                   "test",
		PostgresURL:              "postgres://localhost:5432/test",
		DBMaxOpenConns:           1,
		DBMaxIdleConns:           1,
		DBConnMaxLifetime:        time.Minute,
		JWTSecret:                strings.Repeat("s", 64),
		AccessTokenTTL:           15 * time.Minute,
		RefreshTokenTTL:          7 * 24 * time.Hour,
		RefreshCookieName:        "__Host_refresh",
		LoginRateLimitPerMin:     5,
		LoginRateLimitBurst:      10,
		UploadRateLimitPerMin:    10,
		UploadRateLimitBurst:     10,
		KnowledgeCacheTTL:        time.Minute,
		KnowledgeRateLimitPerMin: 6,
		KnowledgeRateLimitBurst:  30,
		ChatRateLimitPerMin:      6,
		ChatRateLimitBurst:       30,
		ChatModel:                "mock-model",
		Upload: config.UploadConfig{
			MaxBytes:    5 * 1024 * 1024,
			AllowedMIME: []string{"image/png", "image/jpeg", "image/webp", "image/svg+xml"},
			AllowSVG:    false,
		},
		Storage: config.StorageConfig{
			Driver: config.StorageDriverSupabase,
			Supabase: config.SupabaseConfig{
				URL:         "http://localhost",
				Bucket:      "test",
				ServiceRole: "test",
				PublicURL:   "http://localhost/storage/v1/object/public/test",
			},
		},
	}

	srv, err := server.New(dbx, cfg)
	require.NoError(t, err)

	tokenService, err := appauth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}
	return srv, tokenService, cleanup
}
