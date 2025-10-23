package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler responds with a simple heartbeat payload.
type HealthHandler struct{}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandleHealth writes a JSON payload to indicate the API is alive.
func (h *HealthHandler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
