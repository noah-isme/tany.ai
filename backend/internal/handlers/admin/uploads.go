package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
)

// UploadsHandler returns stub responses until storage integration is implemented.
type UploadsHandler struct{}

// NewUploadsHandler constructs UploadsHandler.
func NewUploadsHandler() *UploadsHandler {
	return &UploadsHandler{}
}

// Create responds with stub message for future upload integration.
func (h *UploadsHandler) Create(c *gin.Context) {
	httpapi.RespondData(c, http.StatusOK, gin.H{
		"message": "upload stub",
	})
}
