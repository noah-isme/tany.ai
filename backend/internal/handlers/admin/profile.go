package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// ProfileHandler handles admin profile endpoints.
type ProfileHandler struct {
	repo       repos.ProfileRepository
	invalidate func()
}

// NewProfileHandler constructs a ProfileHandler.
func NewProfileHandler(repo repos.ProfileRepository, invalidate func()) *ProfileHandler {
	ensureValidators()
	return &ProfileHandler{repo: repo, invalidate: invalidate}
}

// Get returns the current profile or 404 if missing.
func (h *ProfileHandler) Get(c *gin.Context) {
	profile, err := h.repo.Get(c.Request.Context())
	if handleRepoError(c, err) {
		return
	}
	httpapi.RespondData(c, http.StatusOK, dto.NewProfileResponse(profile))
}

// Put creates or updates the profile.
func (h *ProfileHandler) Put(c *gin.Context) {
	var req dto.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	model := req.ToModel(uuid.Nil)
	profile, err := h.repo.Upsert(c.Request.Context(), model)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewProfileResponse(profile))
}
