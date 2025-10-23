package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ProfileRequest defines the payload for creating or updating a profile.
type ProfileRequest struct {
	ID        *uuid.UUID `json:"id"`
	Name      string     `json:"name" binding:"required,min=2,max=120"`
	Title     string     `json:"title" binding:"required,min=2,max=160"`
	Bio       string     `json:"bio" binding:"omitempty,max=2000"`
	Email     string     `json:"email" binding:"omitempty,email"`
	Phone     string     `json:"phone" binding:"omitempty,max=64"`
	Location  string     `json:"location" binding:"omitempty,max=160"`
	AvatarURL string     `json:"avatar_url" binding:"omitempty,url"`
}

// ProfileResponse represents the response payload for profile endpoints.
type ProfileResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Bio       string    `json:"bio"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Location  string    `json:"location"`
	AvatarURL string    `json:"avatar_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToModel converts the request into a models.Profile instance.
func (r ProfileRequest) ToModel(existingID uuid.UUID) models.Profile {
	profileID := existingID
	if r.ID != nil {
		profileID = *r.ID
	}

	return models.Profile{
		ID:        profileID,
		Name:      r.Name,
		Title:     r.Title,
		Bio:       r.Bio,
		Email:     r.Email,
		Phone:     r.Phone,
		Location:  r.Location,
		AvatarURL: r.AvatarURL,
	}
}

// NewProfileResponse constructs a ProfileResponse from a model.
func NewProfileResponse(profile models.Profile) ProfileResponse {
	return ProfileResponse{
		ID:        profile.ID.String(),
		Name:      profile.Name,
		Title:     profile.Title,
		Bio:       profile.Bio,
		Email:     profile.Email,
		Phone:     profile.Phone,
		Location:  profile.Location,
		AvatarURL: profile.AvatarURL,
		UpdatedAt: profile.UpdatedAt,
	}
}
