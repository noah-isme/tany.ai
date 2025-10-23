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

// ServiceHandler manages admin service endpoints.
type ServiceHandler struct {
	repo repos.ServiceRepository
}

// NewServiceHandler creates a ServiceHandler.
func NewServiceHandler(repo repos.ServiceRepository) *ServiceHandler {
	ensureValidators()
	return &ServiceHandler{repo: repo}
}

// List returns paginated services.
func (h *ServiceHandler) List(c *gin.Context) {
	params := parseListParams(c)
	services, total, err := h.repo.List(c.Request.Context(), params)
	if handleListError(c, err) {
		return
	}

	responses := make([]dto.ServiceResponse, len(services))
	for i, service := range services {
		responses[i] = dto.NewServiceResponse(service)
	}

	httpapi.RespondList(c, http.StatusOK, responses, params.Page, params.Limit, total)
}

// Create adds a new service entry.
func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	if err := validateServicePrices(req); err != nil {
		respondValidationError(c, err)
		return
	}

	base := models.Service{IsActive: true}
	if req.IsActive != nil {
		base.IsActive = *req.IsActive
	}
	if req.Order != nil {
		base.Order = *req.Order
	}

	model := req.ToModel(uuid.Nil, base)
	created, err := h.repo.Create(c.Request.Context(), model)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	httpapi.RespondData(c, http.StatusCreated, dto.NewServiceResponse(created))
}

// Update modifies an existing service.
func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var req dto.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, err)
		return
	}

	if err := validateServicePrices(req); err != nil {
		respondValidationError(c, err)
		return
	}

	base := models.Service{ID: id, IsActive: true}
	model := req.ToModel(id, base)
	updated, err := h.repo.Update(c.Request.Context(), model)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewServiceResponse(updated))
}

// Delete removes a service.
func (h *ServiceHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		handleRepoError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Reorder updates service ordering.
func (h *ServiceHandler) Reorder(c *gin.Context) {
	var payload []dto.ServiceReorderItem
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondValidationError(c, err)
		return
	}

	updates := make([]models.Service, len(payload))
	for i, item := range payload {
		updates[i] = models.Service{ID: item.ID, Order: item.Order}
	}

	if err := h.repo.Reorder(c.Request.Context(), updates); err != nil {
		handleRepoError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Toggle flips or sets service active state.
func (h *ServiceHandler) Toggle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondValidationError(c, err)
		return
	}

	var desired *bool
	if c.Request.ContentLength > 0 {
		var req dto.ServiceToggleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			respondValidationError(c, err)
			return
		}
		desired = req.IsActive
	}

	service, err := h.repo.Toggle(c.Request.Context(), id, desired)
	if err != nil {
		handleRepoError(c, err)
		return
	}

	httpapi.RespondData(c, http.StatusOK, dto.NewServiceResponse(service))
}

func validateServicePrices(req dto.ServiceRequest) error {
	if req.PriceMin != nil && req.PriceMax != nil && *req.PriceMax < *req.PriceMin {
		return validatorErr("price_max", "must be greater than or equal to price_min")
	}
	return nil
}

func validatorErr(field, message string) error {
	return &fieldError{field: field, message: message}
}

type fieldError struct {
	field   string
	message string
}

func (e *fieldError) Error() string {
	return e.message
}

func (e *fieldError) Field() string {
	return e.field
}
