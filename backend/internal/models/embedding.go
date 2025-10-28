package models

import (
	"time"

	"github.com/google/uuid"
)

// Embedding represents a single semantic vector persisted in storage.
type Embedding struct {
	ID        uuid.UUID  `db:"id"`
	Kind      string     `db:"kind"`
	RefID     *uuid.UUID `db:"ref_id"`
	Content   string     `db:"content"`
	Vector    []float32  `db:"vector"`
	Metadata  JSONB      `db:"metadata"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

// EmbeddingMatch captures similarity search results.
type EmbeddingMatch struct {
	ID       uuid.UUID  `db:"id"`
	Kind     string     `db:"kind"`
	RefID    *uuid.UUID `db:"ref_id"`
	Content  string     `db:"content"`
	Metadata JSONB      `db:"metadata"`
	Score    float64    `db:"score"`
}

// EmbeddingConfig stores personalization settings persisted in the database.
type EmbeddingConfig struct {
	Weight          float64   `json:"weight"`
	LastReindexedAt time.Time `json:"lastReindexedAt"`
	LastResetAt     time.Time `json:"lastResetAt"`
}
