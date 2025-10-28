package models

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsEvent captures granular interaction data for observability dashboards.
type AnalyticsEvent struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	EventType string    `db:"event_type" json:"eventType"`
	Source    string    `db:"source" json:"source"`
	Provider  string    `db:"provider" json:"provider"`
	Duration  int       `db:"duration_ms" json:"durationMs"`
	Success   bool      `db:"success" json:"success"`
	UserAgent string    `db:"user_agent" json:"userAgent"`
	Metadata  JSONB     `db:"metadata" json:"metadata"`
}

// AnalyticsSummary stores aggregated KPI metrics for daily snapshots.
type AnalyticsSummary struct {
	ID                uuid.UUID `db:"id" json:"id"`
	Date              time.Time `db:"date" json:"date"`
	TotalChats        int       `db:"total_chats" json:"totalChats"`
	AvgResponseTimeMS float64   `db:"avg_response_time" json:"avgResponseTime"`
	SuccessRate       float64   `db:"success_rate" json:"successRate"`
	UniqueUsers       int       `db:"unique_users" json:"uniqueUsers"`
	Conversions       int       `db:"conversions" json:"conversions"`
	ProviderBreakdown JSONB     `db:"provider_breakdown" json:"providerBreakdown"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}
