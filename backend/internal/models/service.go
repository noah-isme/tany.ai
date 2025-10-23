package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// Service represents a freelance service offering.
type Service struct {
	ID            uuid.UUID       `db:"id"`
	Name          string          `db:"name"`
	Description   sql.NullString  `db:"description"`
	PriceMin      sql.NullFloat64 `db:"price_min"`
	PriceMax      sql.NullFloat64 `db:"price_max"`
	Currency      sql.NullString  `db:"currency"`
	DurationLabel sql.NullString  `db:"duration_label"`
	IsActive      bool            `db:"is_active"`
	Order         int             `db:"order"`
}
