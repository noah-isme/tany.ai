package models

import "github.com/google/uuid"

// Skill represents a single skill entry manageable via the admin API.
type Skill struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Order int       `db:"order"`
}
