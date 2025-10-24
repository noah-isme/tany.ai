package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	appauth "github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/config"
	admin "github.com/tanydotai/tanyai/backend/internal/handlers/admin"
	"github.com/tanydotai/tanyai/backend/internal/middleware"
	"github.com/tanydotai/tanyai/backend/internal/storage"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type storageStub struct {
	url        string
	err        error
	lastKey    string
	lastType   string
	lastObject []byte
}

func (s *storageStub) Put(_ context.Context, key string, content []byte, contentType string) (string, error) {
	s.lastKey = key
	s.lastType = contentType
	s.lastObject = append([]byte(nil), content...)
	if s.err != nil {
		return "", s.err
	}
	if s.url != "" {
		return s.url, nil
	}
	return "https://cdn.example.com/" + key, nil
}

func TestAdminUploadsUnauthorized(t *testing.T) {
	router, _ := setupUploadRouter(t, &storageStub{}, defaultUploadPolicy(), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminUploadsForbidden(t *testing.T) {
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), nil)

	token, err := tokens.GenerateAccessToken(appauth.Subject{ID: uuidFromString(t, "11111111-1111-1111-1111-111111111111"), Email: "user@example.com", Roles: []string{"editor"}})
	require.NoError(t, err)

	body, boundary := multipartBody(t, "file", "test.png", pngBytes(), "image/png")
	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAdminUploadsTooLarge(t *testing.T) {
	policy := defaultUploadPolicy()
	policy.MaxBytes = 10
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, policy, nil)

	token := mustAdminToken(t, tokens)
	content := bytes.Repeat([]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}, 4)
	body, boundary := multipartBody(t, "file", "big.png", content, "image/png")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}

func TestAdminUploadsUnsupportedType(t *testing.T) {
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), nil)

	token := mustAdminToken(t, tokens)
	body, boundary := multipartBody(t, "file", "test.txt", []byte("hello"), "text/plain")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestAdminUploadsMissingFile(t *testing.T) {
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), nil)

	token := mustAdminToken(t, tokens)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminUploadsSuccess(t *testing.T) {
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), nil)

	token := mustAdminToken(t, tokens)
	body, boundary := multipartBody(t, "file", "avatar.png", pngBytes(), "image/png")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "image/png", store.lastType)
	require.NotEmpty(t, store.lastKey)
	require.NotEmpty(t, store.lastObject)
}

func TestAdminUploadsSvgSanitized(t *testing.T) {
	policy := defaultUploadPolicy()
	policy.AllowSVG = true
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, policy, nil)

	token := mustAdminToken(t, tokens)
	svg := []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"><script>alert('x')</script><rect width="10" height="10" onclick="hack()" href="http://evil"/></svg>`)
	body, boundary := multipartBody(t, "file", "icon.svg", svg, "image/svg+xml")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "image/svg+xml", store.lastType)
	require.NotContains(t, string(store.lastObject), "script")
	require.NotContains(t, string(store.lastObject), "onclick")
	require.NotContains(t, string(store.lastObject), "http://evil")
}

func TestAdminUploadsSvgDeniedByDefault(t *testing.T) {
	store := &storageStub{}
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), nil)

	token := mustAdminToken(t, tokens)
	body, boundary := multipartBody(t, "file", "icon.svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`), "image/svg+xml")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestAdminUploadsRateLimited(t *testing.T) {
	store := &storageStub{}
	limiter := appauth.NewRateLimiter(1, 1, time.Minute)
	router, tokens := setupUploadRouter(t, store, defaultUploadPolicy(), limiter)

	token := mustAdminToken(t, tokens)
	body1, boundary := multipartBody(t, "file", "avatar.png", pngBytes(), "image/png")

	req := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body1)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// second request should be rate limited
	body2, _ := multipartBody(t, "file", "avatar.png", pngBytes(), "image/png")
	req2 := httptest.NewRequest(http.MethodPost, "/api/admin/uploads", body2)
	req2.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	require.Equal(t, http.StatusTooManyRequests, rec2.Code)
}

func setupUploadRouter(t *testing.T, store storage.ObjectStorage, policy config.UploadConfig, limiter *appauth.RateLimiter) (*gin.Engine, *appauth.TokenService) {
	t.Helper()

	tokenService, err := appauth.NewTokenService("this_is_a_super_secret_for_tests_1234567890", time.Hour, time.Hour)
	require.NoError(t, err)

	handler := admin.NewUploadsHandler(store, policy, log.New(io.Discard, "", 0))

	router := gin.New()
	group := router.Group("/api/admin", middleware.Authn(tokenService), middleware.AuthzAdmin())
	group.POST("/uploads", middleware.RateLimitByIP(limiter), handler.Create)
	return router, tokenService
}

func defaultUploadPolicy() config.UploadConfig {
	return config.UploadConfig{
		MaxBytes:    5 * 1024 * 1024,
		AllowedMIME: []string{"image/png", "image/jpeg", "image/webp", "image/svg+xml"},
		AllowSVG:    false,
	}
}

func multipartBody(t *testing.T, field, filename string, content []byte, contentType ...string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	var part io.Writer
	var err error
	if len(contentType) > 0 {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="`+field+`"; filename="`+filename+`"`)
		header.Set("Content-Type", contentType[0])
		part, err = writer.CreatePart(header)
	} else {
		part, err = writer.CreateFormFile(field, filename)
	}
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.Boundary()
}

func pngBytes() []byte {
	return append([]byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xde, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xd7, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xe2, 0x26, 0x05, 0x9b, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
		0x42, 0x60, 0x82,
	}, []byte{}...)
}

func mustAdminToken(t *testing.T, tokens *appauth.TokenService) string {
	t.Helper()
	token, err := tokens.GenerateAccessToken(appauth.Subject{ID: uuidFromString(t, "11111111-1111-1111-1111-111111111111"), Email: "admin@example.com", Roles: []string{"admin"}})
	require.NoError(t, err)
	return token
}

func uuidFromString(t *testing.T, value string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(value)
	require.NoError(t, err)
	return id
}
