package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

// ExternalItemResponse represents external knowledge entries returned to the admin UI.
type ExternalItemResponse struct {
	ID          uuid.UUID      `json:"id"`
	SourceName  string         `json:"sourceName"`
	Kind        string         `json:"kind"`
	Title       string         `json:"title"`
	Summary     string         `json:"summary,omitempty"`
	URL         string         `json:"url"`
	Visible     bool           `json:"visible"`
	PublishedAt *time.Time     `json:"publishedAt,omitempty"`
	Metadata    map[string]any `json:"metadata"`
}

// NewExternalItemResponse converts repository rows into DTOs.
func NewExternalItemResponse(row repos.ExternalItemWithSource) ExternalItemResponse {
	var summary string
	if row.Summary != nil {
		summary = *row.Summary
	}
	if summary == "" && row.Content != nil {
		summary = *row.Content
	}
	var published *time.Time
	if row.PublishedAt != nil {
		published = row.PublishedAt
	}
	metadata := make(map[string]any, len(row.Metadata))
	if row.Metadata != nil {
		for k, v := range row.Metadata {
			metadata[k] = v
		}
	}
	return ExternalItemResponse{
		ID:          row.ID,
		SourceName:  row.SourceName,
		Kind:        row.Kind,
		Title:       row.Title,
		Summary:     summary,
		URL:         row.URL,
		Visible:     row.Visible,
		PublishedAt: published,
		Metadata:    metadata,
	}
}

// ExternalSourceResponse represents an external source configuration.
type ExternalSourceResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	BaseURL      string     `json:"baseUrl"`
	SourceType   string     `json:"sourceType"`
	Enabled      bool       `json:"enabled"`
	LastSyncedAt *time.Time `json:"lastSyncedAt,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
}

// NewExternalSourceResponse converts a model into response format.
func NewExternalSourceResponse(model models.ExternalSource) ExternalSourceResponse {
	return ExternalSourceResponse{
		ID:           model.ID,
		Name:         model.Name,
		BaseURL:      model.BaseURL,
		SourceType:   model.SourceType,
		Enabled:      model.Enabled,
		LastSyncedAt: model.LastSyncedAt,
		LastModified: model.LastModified,
	}
}
