package models

import (
	"time"

	"github.com/google/uuid"
)

// ExternalItem represents normalized structured content fetched from an external source.
type ExternalItem struct {
	ID          uuid.UUID  `db:"id"`
	SourceID    uuid.UUID  `db:"source_id"`
	Kind        string     `db:"kind"`
	Title       string     `db:"title"`
	URL         string     `db:"url"`
	Summary     *string    `db:"summary"`
	Content     *string    `db:"content"`
	Metadata    JSONB      `db:"metadata"`
	PublishedAt *time.Time `db:"published_at"`
	Hash        string     `db:"hash"`
	Visible     bool       `db:"visible"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
