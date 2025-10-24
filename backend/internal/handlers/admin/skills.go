package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/dto"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// SkillHandler manages admin skill endpoints.
type SkillHandler struct {
	repo       repos.SkillRepository
	invalidate func()
}

// NewSkillHandler creates a new SkillHandler.
func NewSkillHandler(repo repos.SkillRepository, invalidate func()) *SkillHandler {
	ensureValidators()
	return &SkillHandler{repo: repo, invalidate: invalidate}
}

// List returns paginated list of skills.
func (h *SkillHandler) List(c *gin.Context) {
	params := parseListParams(c)
	skills, total, err := h.repo.List(c.Request.Context(), params)
	if handleListError(c, err) {
		return
	}

	responses := make([]dto.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = dto.NewSkillResponse(skill)
	}

	httpapi.RespondList(c, http.StatusOK, responses, params.Page, params.Limit, total)
}

// Create adds a new skill entry.
func (h *SkillHandler) Create(c *gin.Context) {
	var req dto.SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	skill := models.Skill{Name: req.Name}
	if req.Order != nil {
		skill.Order = *req.Order
	}

	created, err := h.repo.Create(c.Request.Context(), skill)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusCreated, dto.NewSkillResponse(created))
}

// Update modifies an existing skill.
func (h *SkillHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var req dto.SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	skill := models.Skill{ID: id, Name: req.Name}
	if req.Order != nil {
		skill.Order = *req.Order
	}

	updated, err := h.repo.Update(c.Request.Context(), skill)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewSkillResponse(updated))
}

// Delete removes a skill.
func (h *SkillHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	c.Status(http.StatusNoContent)
}

// Reorder updates ordering of skills.
func (h *SkillHandler) Reorder(c *gin.Context) {
	var payload []dto.SkillReorderItem
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondValidationError(c, err)
		return
	}

	updates := make([]models.Skill, len(payload))
	for i, item := range payload {
		updates[i] = models.Skill{ID: item.ID, Order: item.Order}
	}

	if err := h.repo.Reorder(c.Request.Context(), updates); err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	c.Status(http.StatusNoContent)
}
