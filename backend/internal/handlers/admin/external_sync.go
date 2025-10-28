package admin

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/ingest"
)

// IngestService abstracts external sync behaviour for easier testing.
type IngestService interface {
	Sync(ctx context.Context, source ingest.Source) (ingest.Result, error)
}

// ExternalSourceHandler manages admin endpoints for external sources.
type ExternalSourceHandler struct {
	sources    repos.ExternalSourceRepository
	items      repos.ExternalItemRepository
	ingest     IngestService
	invalidate func()
}

// NewExternalSourceHandler constructs a handler.
func NewExternalSourceHandler(sources repos.ExternalSourceRepository, items repos.ExternalItemRepository, svc IngestService, invalidate func()) *ExternalSourceHandler {
	return &ExternalSourceHandler{sources: sources, items: items, ingest: svc, invalidate: invalidate}
}

// List returns configured external sources.
func (h *ExternalSourceHandler) List(c *gin.Context) {
	params := parseListParams(c)
	rows, total, err := h.sources.List(c.Request.Context(), params)
	if handleListError(c, err) {
		return
	}

	responses := make([]dto.ExternalSourceResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, dto.NewExternalSourceResponse(row))
	}

	httpapi.RespondList(c, http.StatusOK, responses, params.Page, params.Limit, total)
}

// Sync triggers ingestion for a specific source.
func (h *ExternalSourceHandler) Sync(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	source, err := h.sources.Get(c.Request.Context(), id)
	if err != nil {
		handleRepoError(c, err)
		return
	}
	if !source.Enabled {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "source disabled", nil)
		return
	}

	parsedURL, err := parseBaseURL(source.BaseURL)
	if err != nil {
		respondValidationError(c, err)
		return
	}

	ingestSource := ingest.Source{
		ID:           source.ID,
		Name:         source.Name,
		BaseURL:      parsedURL,
		SourceType:   source.SourceType,
		ETag:         source.ETag,
		LastModified: source.LastModified,
	}

	result, err := h.ingest.Sync(c.Request.Context(), ingestSource)
	if err != nil {
		if errors.Is(err, ingest.ErrNotModified) {
			httpapi.RespondData(c, http.StatusOK, gin.H{
				"message":       "no changes",
				"itemsUpserted": 0,
			})
			return
		}
		httpapi.RespondError(c, http.StatusBadGateway, httpapi.ErrorCodeExternal, "failed to sync source", nil)
		return
	}

	if len(result.Items) > 0 {
		if err := h.items.Upsert(c.Request.Context(), result.Items); err != nil {
			handleRepoError(c, err)
			return
		}
	}

	if err := h.sources.UpdateSyncState(c.Request.Context(), source.ID, result.ETag, result.LastModified, result.FetchedAt); err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusOK, gin.H{
		"message":       "sync completed",
		"itemsUpserted": len(result.Items),
		"etag":          result.ETag,
		"lastModified":  result.LastModified,
	})
}

func parseBaseURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("base url required")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed, nil
}

// SeedDefaults ensures default sources exist using repository implementation.
func SeedDefaults(repo repos.ExternalSourceRepository, defaults []models.ExternalSource) func() error {
	return func() error {
		return repo.EnsureDefaults(context.Background(), defaults)
	}
}
