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

// ProjectHandler manages admin project endpoints.
type ProjectHandler struct {
	repo       repos.ProjectRepository
	invalidate func()
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(repo repos.ProjectRepository, invalidate func()) *ProjectHandler {
	ensureValidators()
	return &ProjectHandler{repo: repo, invalidate: invalidate}
}

// List returns paginated projects.
func (h *ProjectHandler) List(c *gin.Context) {
	params := parseListParams(c)
	projects, total, err := h.repo.List(c.Request.Context(), params)
	if handleListError(c, err) {
		return
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = dto.NewProjectResponse(project)
	}

	httpapi.RespondList(c, http.StatusOK, responses, params.Page, params.Limit, total)
}

// Create adds a new project entry.
func (h *ProjectHandler) Create(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	model := req.ToModel(uuid.Nil, models.Project{})
	created, err := h.repo.Create(c.Request.Context(), model)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusCreated, dto.NewProjectResponse(created))
}

// Update modifies a project entry.
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	model := req.ToModel(id, models.Project{ID: id})
	updated, err := h.repo.Update(c.Request.Context(), model)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewProjectResponse(updated))
}

// Delete removes a project entry.
func (h *ProjectHandler) Delete(c *gin.Context) {
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

// Reorder updates ordering of projects.
func (h *ProjectHandler) Reorder(c *gin.Context) {
	var payload []dto.ProjectReorderItem
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondValidationError(c, err)
		return
	}

	updates := make([]models.Project, len(payload))
	for i, item := range payload {
		updates[i] = models.Project{ID: item.ID, Order: item.Order}
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

// Feature sets the featured status of a project.
func (h *ProjectHandler) Feature(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var req dto.ProjectFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	project, err := h.repo.SetFeatured(c.Request.Context(), id, req.IsFeatured)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	if h.invalidate != nil {
		h.invalidate()
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewProjectResponse(project))
}
