package models

import (
	"time"

	"github.com/google/uuid"
)

// Profile represents the persisted profile entity for the admin API.
type Profile struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Title     string    `db:"title"`
	Bio       string    `db:"bio"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Location  string    `db:"location"`
	AvatarURL string    `db:"avatar_url"`
	UpdatedAt time.Time `db:"updated_at"`
}
