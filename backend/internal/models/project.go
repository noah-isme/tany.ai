package models

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Project represents a project portfolio entry.
type Project struct {
	ID          uuid.UUID      `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	TechStack   pq.StringArray `db:"tech_stack"`
	ImageURL    sql.NullString `db:"image_url"`
	ProjectURL  sql.NullString `db:"project_url"`
	Category    sql.NullString `db:"category"`
	Order       int            `db:"order"`
	IsFeatured  bool           `db:"is_featured"`
}
