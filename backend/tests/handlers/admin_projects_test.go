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

func TestAdminProjectsFeatureSuccess(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &projectRepoStub{
		featureResult: models.Project{ID: uuid.New(), Title: "Project", IsFeatured: true},
	}
	router := setupProjectsRouter(repo)

	payload := map[string]bool{"is_featured": true}
	body, _ := json.Marshal(payload)
	url := "/api/admin/projects/" + repo.featureResult.ID.String() + "/feature"
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminProjectsCreateValidationError(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &projectRepoStub{}
	router := setupProjectsRouter(repo)

	payload := map[string]string{"title": ""}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminProjectsFeatureNotFound(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &projectRepoStub{featureErr: repos.ErrNotFound}
	router := setupProjectsRouter(repo)

	payload := map[string]bool{"is_featured": true}
	body, _ := json.Marshal(payload)
	url := "/api/admin/projects/" + uuid.New().String() + "/feature"
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminProjectsGuardUnauthorized(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "true")
	t.Setenv("ADMIN_GUARD_MODE", "")
	router := setupProjectsRouter(&projectRepoStub{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminProjectsGuardForbidden(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "true")
	t.Setenv("ADMIN_GUARD_MODE", "forbidden")
	router := setupProjectsRouter(&projectRepoStub{})

	req := httptest.NewRequest(http.MethodGet, "/api/admin/projects", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func setupProjectsRouter(repo repos.ProjectRepository) *gin.Engine {
	handler := admin.NewProjectHandler(repo)
	router := gin.New()
	group := router.Group("/api/admin", middleware.AuthzAdminStub())
	projects := group.Group("/projects")
	projects.GET("", handler.List)
	projects.POST("", handler.Create)
	projects.PUT(":id", handler.Update)
	projects.DELETE(":id", handler.Delete)
	projects.PATCH("/reorder", handler.Reorder)
	projects.PATCH(":id/feature", handler.Feature)
	return router
}

type projectRepoStub struct {
	listProjects  []models.Project
	listTotal     int64
	listErr       error
	createResult  models.Project
	createErr     error
	updateResult  models.Project
	updateErr     error
	deleteErr     error
	reorderErr    error
	featureResult models.Project
	featureErr    error
}

func (p *projectRepoStub) List(ctx context.Context, params repos.ListParams) ([]models.Project, int64, error) {
	return p.listProjects, p.listTotal, p.listErr
}

func (p *projectRepoStub) Create(ctx context.Context, project models.Project) (models.Project, error) {
	if p.createErr != nil {
		return models.Project{}, p.createErr
	}
	if p.createResult.ID == uuid.Nil {
		project.ID = uuid.New()
		return project, nil
	}
	return p.createResult, nil
}

func (p *projectRepoStub) Update(ctx context.Context, project models.Project) (models.Project, error) {
	if p.updateErr != nil {
		return models.Project{}, p.updateErr
	}
	if p.updateResult.ID == uuid.Nil {
		return project, nil
	}
	return p.updateResult, nil
}

func (p *projectRepoStub) Delete(ctx context.Context, id uuid.UUID) error {
	return p.deleteErr
}

func (p *projectRepoStub) Reorder(ctx context.Context, pairs []models.Project) error {
	return p.reorderErr
}

func (p *projectRepoStub) SetFeatured(ctx context.Context, id uuid.UUID, featured bool) (models.Project, error) {
	if p.featureErr != nil {
		return models.Project{}, p.featureErr
	}
	if p.featureResult.ID == uuid.Nil {
		return models.Project{ID: id, Title: "stub", IsFeatured: featured}, nil
	}
	return p.featureResult, nil
}
