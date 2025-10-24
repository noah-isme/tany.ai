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
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

func TestAdminServicesToggleSuccess(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &serviceRepoStub{
		toggleResult: models.Service{ID: uuid.New(), Name: "Service", IsActive: true},
	}
	router := setupServicesRouter(repo)

	url := "/api/admin/services/" + repo.toggleResult.ID.String() + "/toggle"
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminServicesCreateInvalidPrice(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &serviceRepoStub{}
	router := setupServicesRouter(repo)

	payload := map[string]interface{}{
		"name":      "Invalid Service",
		"price_min": 100,
		"price_max": 50,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/services", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminServicesToggleNotFound(t *testing.T) {
	t.Setenv("ENABLE_ADMIN_GUARD", "false")
	repo := &serviceRepoStub{toggleErr: repos.ErrNotFound}
	router := setupServicesRouter(repo)

	url := "/api/admin/services/" + uuid.New().String() + "/toggle"
	req := httptest.NewRequest(http.MethodPatch, url, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func setupServicesRouter(repo repos.ServiceRepository) *gin.Engine {
	handler := admin.NewServiceHandler(repo, func() {})
	router := gin.New()
	group := router.Group("/api/admin")
	services := group.Group("/services")
	services.GET("", handler.List)
	services.POST("", handler.Create)
	services.PUT(":id", handler.Update)
	services.DELETE(":id", handler.Delete)
	services.PATCH("/reorder", handler.Reorder)
	services.PATCH(":id/toggle", handler.Toggle)
	return router
}

type serviceRepoStub struct {
	listServices  []models.Service
	listTotal     int64
	listErr       error
	createResult  models.Service
	createErr     error
	updateResult  models.Service
	updateErr     error
	deleteErr     error
	reorderErr    error
	toggleResult  models.Service
	toggleErr     error
	toggleDesired *bool
}

func (s *serviceRepoStub) List(ctx context.Context, params repos.ListParams) ([]models.Service, int64, error) {
	return s.listServices, s.listTotal, s.listErr
}

func (s *serviceRepoStub) Create(ctx context.Context, service models.Service) (models.Service, error) {
	if s.createErr != nil {
		return models.Service{}, s.createErr
	}
	if s.createResult.ID == uuid.Nil {
		service.ID = uuid.New()
		return service, nil
	}
	return s.createResult, nil
}

func (s *serviceRepoStub) Update(ctx context.Context, service models.Service) (models.Service, error) {
	if s.updateErr != nil {
		return models.Service{}, s.updateErr
	}
	if s.updateResult.ID == uuid.Nil {
		return service, nil
	}
	return s.updateResult, nil
}

func (s *serviceRepoStub) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteErr
}

func (s *serviceRepoStub) Reorder(ctx context.Context, pairs []models.Service) error {
	return s.reorderErr
}

func (s *serviceRepoStub) Toggle(ctx context.Context, id uuid.UUID, desired *bool) (models.Service, error) {
	s.toggleDesired = desired
	if s.toggleErr != nil {
		return models.Service{}, s.toggleErr
	}
	if s.toggleResult.ID == uuid.Nil {
		return models.Service{ID: id, Name: "stub", IsActive: true}, nil
	}
	return s.toggleResult, nil
}
