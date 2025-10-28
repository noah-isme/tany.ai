package embedding

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

// KnowledgeProvider exposes the subset of Aggregator behaviour required by the handler.
type KnowledgeProvider interface {
	Get(ctx context.Context) (kb.KnowledgeBase, string, bool, error)
}

// Handler exposes admin endpoints for managing personalization embeddings.
type Handler struct {
	service    *Service
	knowledge  KnowledgeProvider
	invalidate func()
}

// NewHandler creates a new personalization handler.
func NewHandler(service *Service, knowledge KnowledgeProvider, invalidate func()) *Handler {
	return &Handler{service: service, knowledge: knowledge, invalidate: invalidate}
}

// Summary returns personalization status metrics.
func (h *Handler) Summary(c *gin.Context) {
	summary, err := h.service.Summary(c.Request.Context())
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load personalization summary", nil)
		return
	}
	httpapi.RespondData(c, http.StatusOK, summary)
}

// Reindex rebuilds embeddings from the current knowledge base.
func (h *Handler) Reindex(c *gin.Context) {
	if h.service == nil || h.knowledge == nil {
		httpapi.RespondError(c, http.StatusServiceUnavailable, httpapi.ErrorCodeInternal, "personalization service unavailable", nil)
		return
	}
	base, _, _, err := h.knowledge.Get(c.Request.Context())
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load knowledge base", nil)
		return
	}
	count, err := h.service.Reindex(c.Request.Context(), base)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, err.Error(), nil)
		return
	}
	if h.invalidate != nil {
		h.invalidate()
	}
	httpapi.RespondData(c, http.StatusAccepted, gin.H{"indexed": count})
}

// Reset removes all embeddings and clears caches.
func (h *Handler) Reset(c *gin.Context) {
	if h.service == nil {
		httpapi.RespondError(c, http.StatusServiceUnavailable, httpapi.ErrorCodeInternal, "personalization service unavailable", nil)
		return
	}
	if err := h.service.Reset(c.Request.Context()); err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, err.Error(), nil)
		return
	}
	if h.invalidate != nil {
		h.invalidate()
	}
	httpapi.RespondData(c, http.StatusAccepted, gin.H{"status": "reset"})
}

// UpdateWeight adjusts personalization weight multiplier.
func (h *Handler) UpdateWeight(c *gin.Context) {
	if h.service == nil {
		httpapi.RespondError(c, http.StatusServiceUnavailable, httpapi.ErrorCodeInternal, "personalization service unavailable", nil)
		return
	}
	var payload struct {
		Weight float64 `json:"weight" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "weight is required", nil)
		return
	}
	weight, err := h.service.SetWeight(c.Request.Context(), payload.Weight)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, err.Error(), nil)
		return
	}
	httpapi.RespondData(c, http.StatusOK, gin.H{"weight": weight})
}
