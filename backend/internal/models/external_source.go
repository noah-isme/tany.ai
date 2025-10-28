package models

import (
    "time"

    "github.com/google/uuid"
)

// ExternalSource represents a configured external knowledge provider.
type ExternalSource struct {
    ID           uuid.UUID `db:"id"`
    Name         string    `db:"name"`
    BaseURL      string    `db:"base_url"`
    SourceType   string    `db:"source_type"`
    Enabled      bool      `db:"enabled"`
    ETag         *string   `db:"etag"`
    LastModified *time.Time `db:"last_modified"`
    LastSyncedAt *time.Time `db:"last_synced_at"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
}
