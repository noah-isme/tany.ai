package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthnMissingToken(t *testing.T) {
	tokenService, err := appauth.NewTokenService(strings.Repeat("m", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	router := gin.New()
	router.GET("/admin", Authn(tokenService), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthnInvalidToken(t *testing.T) {
	tokenService, err := appauth.NewTokenService(strings.Repeat("n", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	router := gin.New()
	router.GET("/admin", Authn(tokenService), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthnExpiredToken(t *testing.T) {
	tokenService, err := appauth.NewTokenService(strings.Repeat("o", 64), time.Second, time.Hour)
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	base := time.Now()
	tokenService.SetClock(func() time.Time { return base })
	token, err := tokenService.GenerateAccessToken(appauth.Subject{ID: uuid.New(), Email: "user@example.com"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	tokenService.SetClock(func() time.Time { return base.Add(appauth.ClockSkew + 2*time.Second) })

	router := gin.New()
	router.GET("/admin", Authn(tokenService), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthzAdminRejectsNonAdmin(t *testing.T) {
	tokenService, err := appauth.NewTokenService(strings.Repeat("p", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	token, err := tokenService.GenerateAccessToken(appauth.Subject{ID: uuid.New(), Email: "user@example.com", Roles: []string{"user"}})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	router := gin.New()
	router.GET("/admin", Authn(tokenService), AuthzAdmin(), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestAuthzAdminAllowsAdmin(t *testing.T) {
	tokenService, err := appauth.NewTokenService(strings.Repeat("q", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("token service: %v", err)
	}
	token, err := tokenService.GenerateAccessToken(appauth.Subject{ID: uuid.New(), Email: "admin@example.com", Roles: []string{"admin"}})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	router := gin.New()
	router.GET("/admin", Authn(tokenService), AuthzAdmin(), func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
