package dto

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ProjectRequest captures payload for creating or updating a project entry.
type ProjectRequest struct {
	Title       string   `json:"title" binding:"required,min=2,max=160"`
	Description string   `json:"description" binding:"omitempty,max=4000"`
	TechStack   []string `json:"tech_stack" binding:"omitempty,dive,max=32"`
	ImageURL    string   `json:"image_url" binding:"omitempty,url"`
	ProjectURL  string   `json:"project_url" binding:"omitempty,url"`
	Category    string   `json:"category" binding:"omitempty,max=80"`
	Order       *int     `json:"order" binding:"omitempty,min=0"`
	IsFeatured  *bool    `json:"is_featured"`
}

// ProjectResponse is returned by project endpoints.
type ProjectResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
	ImageURL    string   `json:"image_url"`
	ProjectURL  string   `json:"project_url"`
	Category    string   `json:"category"`
	Order       int      `json:"order"`
	IsFeatured  bool     `json:"is_featured"`
}

// ProjectReorderItem describes reorder payload.
type ProjectReorderItem struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Order int       `json:"order" binding:"required,min=0"`
}

// ProjectFeatureRequest toggles featured flag.
type ProjectFeatureRequest struct {
	IsFeatured bool `json:"is_featured" binding:"required"`
}

// ToModel converts request to models.Project.
func (r ProjectRequest) ToModel(id uuid.UUID, existing models.Project) models.Project {
	result := existing
	result.ID = id
	result.Title = r.Title
	result.Description = sql.NullString{String: r.Description, Valid: r.Description != ""}
	if len(r.TechStack) > 0 {
		result.TechStack = pq.StringArray(r.TechStack)
	}
	result.ImageURL = sql.NullString{String: r.ImageURL, Valid: r.ImageURL != ""}
	result.ProjectURL = sql.NullString{String: r.ProjectURL, Valid: r.ProjectURL != ""}
	result.Category = sql.NullString{String: r.Category, Valid: r.Category != ""}
	if r.Order != nil {
		result.Order = *r.Order
	}
	if r.IsFeatured != nil {
		result.IsFeatured = *r.IsFeatured
	}
	return result
}

// NewProjectResponse converts model to response struct.
func NewProjectResponse(project models.Project) ProjectResponse {
	return ProjectResponse{
		ID:          project.ID.String(),
		Title:       project.Title,
		Description: project.Description.String,
		TechStack:   []string(project.TechStack),
		ImageURL:    project.ImageURL.String,
		ProjectURL:  project.ProjectURL.String,
		Category:    project.Category.String,
		Order:       project.Order,
		IsFeatured:  project.IsFeatured,
	}
}
