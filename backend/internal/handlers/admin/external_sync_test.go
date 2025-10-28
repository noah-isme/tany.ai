package admin

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
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/ingest"
)

type stubSourceRepo struct {
	listFn       func(context.Context, repos.ListParams) ([]models.ExternalSource, int64, error)
	getFn        func(context.Context, uuid.UUID) (models.ExternalSource, error)
	updateSyncFn func(context.Context, uuid.UUID, *string, *time.Time, time.Time) error
}

func (s *stubSourceRepo) List(ctx context.Context, params repos.ListParams) ([]models.ExternalSource, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, params)
	}
	return nil, 0, nil
}

func (s *stubSourceRepo) Get(ctx context.Context, id uuid.UUID) (models.ExternalSource, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return models.ExternalSource{}, repos.ErrNotFound
}

func (s *stubSourceRepo) FindByBaseURL(context.Context, string) (models.ExternalSource, error) {
	return models.ExternalSource{}, repos.ErrNotFound
}

func (s *stubSourceRepo) Create(context.Context, models.ExternalSource) (models.ExternalSource, error) {
	return models.ExternalSource{}, repos.ErrNotFound
}

func (s *stubSourceRepo) Update(context.Context, models.ExternalSource) (models.ExternalSource, error) {
	return models.ExternalSource{}, repos.ErrNotFound
}

func (s *stubSourceRepo) UpdateSyncState(ctx context.Context, id uuid.UUID, etag *string, lastModified *time.Time, syncedAt time.Time) error {
	if s.updateSyncFn != nil {
		return s.updateSyncFn(ctx, id, etag, lastModified, syncedAt)
	}
	return nil
}

func (s *stubSourceRepo) SetEnabled(context.Context, uuid.UUID, bool) (models.ExternalSource, error) {
	return models.ExternalSource{}, repos.ErrNotFound
}

func (s *stubSourceRepo) EnsureDefaults(context.Context, []models.ExternalSource) error {
	return nil
}

type stubItemRepo struct {
	listFn       func(context.Context, repos.ExternalItemListParams) ([]repos.ExternalItemWithSource, int64, error)
	upsertFn     func(context.Context, []models.ExternalItem) error
	visibilityFn func(context.Context, uuid.UUID, bool) (repos.ExternalItemWithSource, error)
}

func (s *stubItemRepo) Upsert(ctx context.Context, items []models.ExternalItem) error {
	if s.upsertFn != nil {
		return s.upsertFn(ctx, items)
	}
	return nil
}

func (s *stubItemRepo) List(ctx context.Context, params repos.ExternalItemListParams) ([]repos.ExternalItemWithSource, int64, error) {
	if s.listFn != nil {
		return s.listFn(ctx, params)
	}
	return nil, 0, nil
}

func (s *stubItemRepo) SetVisibility(ctx context.Context, id uuid.UUID, visible bool) (repos.ExternalItemWithSource, error) {
	if s.visibilityFn != nil {
		return s.visibilityFn(ctx, id, visible)
	}
	return repos.ExternalItemWithSource{}, repos.ErrNotFound
}

type stubIngestService struct {
	syncFn func(context.Context, ingest.Source) (ingest.Result, error)
}

func (s *stubIngestService) Sync(ctx context.Context, src ingest.Source) (ingest.Result, error) {
	if s.syncFn != nil {
		return s.syncFn(ctx, src)
	}
	return ingest.Result{}, nil
}

func TestExternalSourceHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srcID := uuid.New()
	repo := &stubSourceRepo{
		listFn: func(context.Context, repos.ListParams) ([]models.ExternalSource, int64, error) {
			now := time.Now()
			return []models.ExternalSource{{
				ID:           srcID,
				Name:         "noahis.me",
				BaseURL:      "https://noahis.me",
				SourceType:   "auto",
				Enabled:      true,
				LastSyncedAt: &now,
			}}, 1, nil
		},
	}

	handler := NewExternalSourceHandler(repo, &stubItemRepo{}, &stubIngestService{}, func() {})

	engine := gin.New()
	engine.GET("/sources", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	var payload struct {
		Items []map[string]any `json:"items"`
		Total int64            `json:"total"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, int64(1), payload.Total)
	require.Len(t, payload.Items, 1)
	require.Equal(t, "noahis.me", payload.Items[0]["name"])
}

func TestExternalSourceHandlerSync(t *testing.T) {
	gin.SetMode(gin.TestMode)

	srcID := uuid.New()
	now := time.Now()
	repo := &stubSourceRepo{
		getFn: func(context.Context, uuid.UUID) (models.ExternalSource, error) {
			return models.ExternalSource{
				ID:         srcID,
				Name:       "noahis.me",
				BaseURL:    "https://noahis.me",
				SourceType: "auto",
				Enabled:    true,
			}, nil
		},
		updateSyncFn: func(context.Context, uuid.UUID, *string, *time.Time, time.Time) error {
			return nil
		},
	}

	itemRepo := &stubItemRepo{
		upsertFn: func(context.Context, []models.ExternalItem) error { return nil },
	}

	ingestSvc := &stubIngestService{
		syncFn: func(context.Context, ingest.Source) (ingest.Result, error) {
			summary := "Ringkasan"
			item := models.ExternalItem{
				SourceID: srcID,
				Kind:     "post",
				Title:    "Hello",
				URL:      "https://noahis.me/post",
				Summary:  &summary,
				Visible:  true,
			}
			return ingest.Result{Items: []models.ExternalItem{item}, FetchedAt: now}, nil
		},
	}

	invalidated := false
	handler := NewExternalSourceHandler(repo, itemRepo, ingestSvc, func() { invalidated = true })

	engine := gin.New()
	engine.POST("/sources/:id/sync", handler.Sync)

	req := httptest.NewRequest(http.MethodPost, "/sources/"+srcID.String()+"/sync", nil)
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.True(t, invalidated)
	var payload struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, float64(1), payload.Data["itemsUpserted"])
	require.Equal(t, "sync completed", payload.Data["message"])
}

func TestExternalItemHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemID := uuid.New()
	repo := &stubItemRepo{
		listFn: func(context.Context, repos.ExternalItemListParams) ([]repos.ExternalItemWithSource, int64, error) {
			summary := "Summary"
			return []repos.ExternalItemWithSource{{
				ExternalItem: models.ExternalItem{
					ID:      itemID,
					Title:   "Post",
					URL:     "https://noahis.me/post",
					Summary: &summary,
					Visible: true,
				},
				SourceName:    "noahis.me",
				SourceBaseURL: "https://noahis.me",
			}}, 1, nil
		},
	}

	handler := NewExternalItemHandler(repo, func() {})

	engine := gin.New()
	engine.GET("/items", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	var payload struct {
		Items []map[string]any `json:"items"`
		Total int64            `json:"total"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, int64(1), payload.Total)
	require.Len(t, payload.Items, 1)
	require.Equal(t, "Post", payload.Items[0]["title"])
}

func TestExternalItemHandlerToggleVisibility(t *testing.T) {
	gin.SetMode(gin.TestMode)

	itemID := uuid.New()
	repo := &stubItemRepo{
		visibilityFn: func(context.Context, uuid.UUID, bool) (repos.ExternalItemWithSource, error) {
			return repos.ExternalItemWithSource{
				ExternalItem: models.ExternalItem{
					ID:      itemID,
					Title:   "Post",
					URL:     "https://noahis.me/post",
					Visible: false,
				},
				SourceName: "noahis.me",
			}, nil
		},
	}

	invalidated := false
	handler := NewExternalItemHandler(repo, func() { invalidated = true })

	engine := gin.New()
	engine.PATCH("/items/:id/visibility", handler.ToggleVisibility)

	body, _ := json.Marshal(map[string]bool{"visible": false})
	req := httptest.NewRequest(http.MethodPatch, "/items/"+itemID.String()+"/visibility", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.True(t, invalidated)
	var payload struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, false, payload.Data["visible"])
}
