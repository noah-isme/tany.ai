package dto

import "github.com/tanydotai/tanyai/backend/internal/models"

// UserResponse represents the authenticated user payload.
type UserResponse struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  *string  `json:"name,omitempty"`
	Roles []string `json:"roles"`
}

// NewUserResponse builds a UserResponse from a user model and roles.
func NewUserResponse(user models.User, roles []string) UserResponse {
	return UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Roles: append([]string(nil), roles...),
	}
}
