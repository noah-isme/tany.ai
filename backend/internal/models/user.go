package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user account.
type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Name         *string   `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// RefreshToken represents a persisted refresh token record.
type RefreshToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
}
