package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	admin "github.com/tanydotai/tanyai/backend/internal/handlers/admin"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAdminProfileGetSuccess(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &profileRepoStub{
		profile: models.Profile{
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Name:      "John Doe",
			Title:     "Freelancer",
			Email:     "john@example.com",
			UpdatedAt: time.Now(),
		},
	}

	router := setupProfileRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data := body["data"].(map[string]interface{})
	require.Equal(t, "John Doe", data["name"])
}

func TestAdminProfilePutValidationError(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &profileRepoStub{}
	router := setupProfileRouter(repo)

	payload := map[string]string{"name": "", "title": ""}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminProfilePutSuccess(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &profileRepoStub{
		upsertResult: models.Profile{
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Name:      "Jane Doe",
			Title:     "Designer",
			Email:     "jane@example.com",
			UpdatedAt: time.Now(),
		},
	}
	router := setupProfileRouter(repo)

	payload := map[string]string{
		"name":       "Jane Doe",
		"title":      "Designer",
		"email":      "jane@example.com",
		"avatar_url": "https://example.com/avatar.png",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "Jane Doe", repo.lastUpsert.Name)
}

func TestAdminProfileGetNotFound(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &profileRepoStub{getErr: repos.ErrNotFound}
	router := setupProfileRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func setupProfileRouter(repo repos.ProfileRepository) *gin.Engine {
	handler := admin.NewProfileHandler(repo)
	router := gin.New()
	group := router.Group("/api/admin")
	group.GET("/profile", handler.Get)
	group.PUT("/profile", handler.Put)
	return router
}

type profileRepoStub struct {
	profile      models.Profile
	getErr       error
	upsertResult models.Profile
	upsertErr    error
	lastUpsert   models.Profile
}

func (s *profileRepoStub) Get(ctx context.Context) (models.Profile, error) {
	return s.profile, s.getErr
}

func (s *profileRepoStub) Upsert(ctx context.Context, profile models.Profile) (models.Profile, error) {
	s.lastUpsert = profile
	if s.upsertErr != nil {
		return models.Profile{}, s.upsertErr
	}
	if s.upsertResult.ID == uuid.Nil {
		return profile, nil
	}
	return s.upsertResult, nil
}
