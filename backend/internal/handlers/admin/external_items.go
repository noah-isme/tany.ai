package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// ExternalItemHandler exposes admin endpoints for external content.
type ExternalItemHandler struct {
	repo       repos.ExternalItemRepository
	invalidate func()
}

// NewExternalItemHandler constructs a handler.
func NewExternalItemHandler(repo repos.ExternalItemRepository, invalidate func()) *ExternalItemHandler {
	return &ExternalItemHandler{repo: repo, invalidate: invalidate}
}

// List returns paginated external items with optional filters.
func (h *ExternalItemHandler) List(c *gin.Context) {
	params := parseListParams(c)
	filter := repos.ExternalItemListParams{ListParams: params, Kind: strings.TrimSpace(c.Query("kind")), Search: c.Query("q")}
	if source := strings.TrimSpace(c.Query("sourceId")); source != "" {
		id, err := uuid.Parse(source)
		if err != nil {
			respondValidationError(c, err)
			return
		}
		filter.SourceID = &id
	}
	if visible := strings.TrimSpace(c.Query("visible")); visible != "" {
		parsed, err := strconv.ParseBool(visible)
		if err != nil {
			respondValidationError(c, err)
			return
		}
		filter.Visible = &parsed
	}

	rows, total, err := h.repo.List(c.Request.Context(), filter)
	if handleListError(c, err) {
		return
	}

	responses := make([]dto.ExternalItemResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, dto.NewExternalItemResponse(row))
	}

	httpapi.RespondList(c, http.StatusOK, responses, params.Page, params.Limit, total)
}

// ToggleVisibility updates item visibility status.
func (h *ExternalItemHandler) ToggleVisibility(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var payload struct {
		Visible bool `json:"visible"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondValidationError(c, err)
		return
	}

	item, err := h.repo.SetVisibility(c.Request.Context(), id, payload.Visible)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	response := dto.NewExternalItemResponse(item)
	httpapi.RespondData(c, http.StatusOK, response)
}
