package dto

import (
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// SkillRequest defines payload for creating or updating a skill.
type SkillRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=80"`
	Order *int   `json:"order" binding:"omitempty,min=0"`
}

// SkillResponse is returned for skill endpoints.
type SkillResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// SkillReorderItem represents a request body entry for reorder endpoints.
type SkillReorderItem struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Order int       `json:"order" binding:"required,min=0"`
}

// NewSkillResponse converts a skill model to response.
func NewSkillResponse(skill models.Skill) SkillResponse {
	return SkillResponse{
		ID:    skill.ID.String(),
		Name:  skill.Name,
		Order: skill.Order,
	}
}
