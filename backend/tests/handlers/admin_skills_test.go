package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	admin "github.com/tanydotai/tanyai/backend/internal/handlers/admin"
	"github.com/tanydotai/tanyai/backend/internal/middleware"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

func TestAdminSkillsListSuccess(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &skillRepoStub{
		listSkills: []models.Skill{
			{ID: uuid.New(), Name: "Golang", Order: 1},
			{ID: uuid.New(), Name: "React", Order: 2},
		},
		listTotal: 2,
	}
	router := setupSkillsRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/skills", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	items := body["items"].([]interface{})
	require.Len(t, items, 2)
}

func TestAdminSkillsCreateValidationError(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &skillRepoStub{}
	router := setupSkillsRouter(repo)

	payload := map[string]string{"name": ""}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminSkillsUpdateNotFound(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &skillRepoStub{updateErr: repos.ErrNotFound}
	router := setupSkillsRouter(repo)

	payload := map[string]interface{}{"name": "Updated", "order": 1}
	body, _ := json.Marshal(payload)

	url := "/api/admin/skills/" + uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminSkillsGuardUnauthorized(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "true")
	t.Setenv("ADMIN_GUARD_MODE", "")
	router := setupSkillsRouter(&skillRepoStub{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/skills", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminSkillsGuardForbidden(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "true")
	t.Setenv("ADMIN_GUARD_MODE", "forbidden")
	router := setupSkillsRouter(&skillRepoStub{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/skills", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func setupSkillsRouter(repo repos.SkillRepository) *gin.Engine {
	handler := admin.NewSkillHandler(repo)
	router := gin.New()
	group := router.Group("/api/admin", middleware.AuthzAdminStub())
	skills := group.Group("/skills")
	skills.GET("", handler.List)
	skills.POST("", handler.Create)
	skills.PUT(":id", handler.Update)
	skills.DELETE(":id", handler.Delete)
	skills.PATCH("/reorder", handler.Reorder)
	return router
}

type skillRepoStub struct {
	listSkills      []models.Skill
	listTotal       int64
	listErr         error
	createResult    models.Skill
	createErr       error
	updateResult    models.Skill
	updateErr       error
	deleteErr       error
	reorderErr      error
	capturedCreate  models.Skill
	capturedUpdate  models.Skill
	capturedReorder []models.Skill
}

func (s *skillRepoStub) List(ctx context.Context, params repos.ListParams) ([]models.Skill, int64, error) {
	return s.listSkills, s.listTotal, s.listErr
}

func (s *skillRepoStub) Create(ctx context.Context, skill models.Skill) (models.Skill, error) {
	s.capturedCreate = skill
	if s.createErr != nil {
		return models.Skill{}, s.createErr
	}
	if s.createResult.ID == uuid.Nil {
		skill.ID = uuid.New()
		return skill, nil
	}
	return s.createResult, nil
}

func (s *skillRepoStub) Update(ctx context.Context, skill models.Skill) (models.Skill, error) {
	s.capturedUpdate = skill
	if s.updateErr != nil {
		return models.Skill{}, s.updateErr
	}
	if s.updateResult.ID == uuid.Nil {
		return skill, nil
	}
	return s.updateResult, nil
}

func (s *skillRepoStub) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteErr
}

func (s *skillRepoStub) Reorder(ctx context.Context, pairs []models.Skill) error {
	s.capturedReorder = pairs
	return s.reorderErr
}
