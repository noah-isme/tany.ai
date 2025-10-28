package ingest

import (
    "errors"
    "net/url"
    "time"

    "github.com/google/uuid"
    "github.com/tanydotai/tanyai/backend/internal/models"
)

// ErrNotModified indicates that the remote source has not changed since the last sync.
var ErrNotModified = errors.New("external source not modified")

// Source describes the configuration for an external source sync cycle.
type Source struct {
    ID           uuid.UUID
    Name         string
    BaseURL      *url.URL
    SourceType   string
    ETag         *string
    LastModified *time.Time
}

// Result captures the outcome of a sync run.
type Result struct {
    Items        []models.ExternalItem
    ETag         *string
    LastModified *time.Time
    FetchedAt    time.Time
}
