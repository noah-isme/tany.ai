package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// HealthHandler responds with a simple heartbeat payload and checks database connectivity.
type HealthHandler struct {
	db *sqlx.DB
}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler(db *sqlx.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HandleHealth writes a JSON payload to indicate the API and database status.
func (h *HealthHandler) HandleHealth(c *gin.Context) {
	response := gin.H{
		"status": "ok",
	}

	if h.db != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := h.db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"database": gin.H{
					"status": "error",
					"error":  err.Error(),
				},
			})
			return
		}

		response["database"] = gin.H{"status": "ok"}
	}

	c.JSON(http.StatusOK, response)
}
