package analytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

// Handler exposes admin analytics endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs a Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Summary returns aggregated analytics metrics for the selected period.
func (h *Handler) Summary(c *gin.Context) {
	filter := RangeFilter{}
	if from := c.Query("from"); from != "" {
		parsed, err := time.Parse(time.RFC3339, from)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid from timestamp", nil)
			return
		}
		filter.Start = parsed
	}
	if to := c.Query("to"); to != "" {
		parsed, err := time.Parse(time.RFC3339, to)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid to timestamp", nil)
			return
		}
		filter.End = parsed
	}
	filter.Source = c.Query("source")
	filter.Provider = c.Query("provider")

	summary, err := h.service.Summary(c.Request.Context(), filter)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load analytics summary", nil)
		return
	}
	httpapi.RespondData(c, http.StatusOK, summary)
}

// Events returns paginated analytics events (chat and custom).
func (h *Handler) Events(c *gin.Context) {
	filter := EventFilter{}
	if from := c.Query("from"); from != "" {
		parsed, err := time.Parse(time.RFC3339, from)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid from timestamp", nil)
			return
		}
		filter.Start = parsed
	}
	if to := c.Query("to"); to != "" {
		parsed, err := time.Parse(time.RFC3339, to)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid to timestamp", nil)
			return
		}
		filter.End = parsed
	}
	filter.Source = c.Query("source")
	filter.Provider = c.Query("provider")
	filter.Type = c.Query("type")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	filter.Limit = limit
	if page <= 1 {
		filter.Offset = 0
	} else {
		filter.Offset = (page - 1) * limit
	}

	result, err := h.service.Events(c.Request.Context(), filter)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load analytics events", nil)
		return
	}
	httpapi.RespondList(c, http.StatusOK, result.Items, page, limit, result.Total)
}

// Leads returns analytics events filtered by lead type.
func (h *Handler) Leads(c *gin.Context) {
	filter := EventFilter{Type: "lead"}
	if from := c.Query("from"); from != "" {
		parsed, err := time.Parse(time.RFC3339, from)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid from timestamp", nil)
			return
		}
		filter.Start = parsed
	}
	if to := c.Query("to"); to != "" {
		parsed, err := time.Parse(time.RFC3339, to)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid to timestamp", nil)
			return
		}
		filter.End = parsed
	}
	filter.Source = c.Query("source")
	filter.Provider = c.Query("provider")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	filter.Limit = limit
	if page <= 1 {
		filter.Offset = 0
	} else {
		filter.Offset = (page - 1) * limit
	}

	result, err := h.service.Events(c.Request.Context(), filter)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load leads analytics", nil)
		return
	}
	httpapi.RespondList(c, http.StatusOK, result.Items, page, limit, result.Total)
}
