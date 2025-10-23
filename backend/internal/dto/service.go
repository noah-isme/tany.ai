package dto

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ServiceRequest describes payload for creating or updating a service.
type ServiceRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=120"`
	Description   string   `json:"description" binding:"omitempty,max=2000"`
	PriceMin      *float64 `json:"price_min" binding:"omitempty,min=0"`
	PriceMax      *float64 `json:"price_max" binding:"omitempty,min=0"`
	Currency      string   `json:"currency" binding:"omitempty,currency_code"`
	DurationLabel string   `json:"duration_label" binding:"omitempty,max=80"`
	IsActive      *bool    `json:"is_active"`
	Order         *int     `json:"order" binding:"omitempty,min=0"`
}

// ServiceResponse represents the JSON output of a service entity.
type ServiceResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PriceMin      *float64 `json:"price_min"`
	PriceMax      *float64 `json:"price_max"`
	Currency      string   `json:"currency"`
	DurationLabel string   `json:"duration_label"`
	IsActive      bool     `json:"is_active"`
	Order         int      `json:"order"`
}

// ServiceToggleRequest captures payload for toggle endpoint.
type ServiceToggleRequest struct {
	IsActive *bool `json:"is_active" binding:"omitempty"`
}

// ServiceReorderItem describes reorder payload entries.
type ServiceReorderItem struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Order int       `json:"order" binding:"required,min=0"`
}

// ToModel converts request to models.Service.
func (r ServiceRequest) ToModel(id uuid.UUID, existing models.Service) models.Service {
	result := existing
	result.ID = id
	result.Name = r.Name
	result.Description = sql.NullString{String: r.Description, Valid: r.Description != ""}
	if r.PriceMin != nil {
		result.PriceMin = sql.NullFloat64{Float64: *r.PriceMin, Valid: true}
	}
	if r.PriceMax != nil {
		result.PriceMax = sql.NullFloat64{Float64: *r.PriceMax, Valid: true}
	}
	if r.Currency != "" {
		result.Currency = sql.NullString{String: strings.ToUpper(r.Currency), Valid: true}
	}
	result.DurationLabel = sql.NullString{String: r.DurationLabel, Valid: r.DurationLabel != ""}
	if r.IsActive != nil {
		result.IsActive = *r.IsActive
	}
	if r.Order != nil {
		result.Order = *r.Order
	}
	return result
}

// NewServiceResponse builds response from model.
func NewServiceResponse(service models.Service) ServiceResponse {
	return ServiceResponse{
		ID:            service.ID.String(),
		Name:          service.Name,
		Description:   service.Description.String,
		PriceMin:      nullFloat64Pointer(service.PriceMin),
		PriceMax:      nullFloat64Pointer(service.PriceMax),
		Currency:      service.Currency.String,
		DurationLabel: service.DurationLabel.String,
		IsActive:      service.IsActive,
		Order:         service.Order,
	}
}

func nullFloat64Pointer(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	value := v.Float64
	return &value
}
