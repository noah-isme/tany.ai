package analytics

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ErrAnalyticsDisabled is returned when analytics module is disabled.
var ErrAnalyticsDisabled = errors.New("analytics disabled")

// RecordChatInput represents telemetry captured for chat requests.
type RecordChatInput struct {
	Timestamp time.Time
	Source    string
	Provider  string
	Duration  time.Duration
	Success   bool
	UserAgent string
	ChatID    uuid.UUID
	Metadata  models.JSONB
}

// RecordEventInput supports recording arbitrary analytics events.
type RecordEventInput struct {
	Timestamp time.Time
	Type      string
	Source    string
	Provider  string
	Duration  time.Duration
	Success   bool
	UserAgent string
	Metadata  models.JSONB
}

// SummaryRange captures aggregated metrics for dashboards.
type SummaryRange struct {
	RangeStart        time.Time                   `json:"rangeStart"`
	RangeEnd          time.Time                   `json:"rangeEnd"`
	TotalChats        int                         `json:"totalChats"`
	AvgResponseTimeMS float64                     `json:"avgResponseTime"`
	SuccessRate       float64                     `json:"successRate"`
	UniqueUsers       int                         `json:"uniqueUsers"`
	Conversions       int                         `json:"conversions"`
	ProviderBreakdown map[string]ProviderSnapshot `json:"providerBreakdown"`
	Daily             []DailySnapshot             `json:"daily"`
}

// ProviderSnapshot summarises provider level metrics.
type ProviderSnapshot struct {
	TotalChats        int     `json:"totalChats"`
	AvgResponseTimeMS float64 `json:"avgResponseTime"`
	SuccessRate       float64 `json:"successRate"`
}

// DailySnapshot summarises metrics per day for charting.
type DailySnapshot struct {
	Date              time.Time `json:"date"`
	TotalChats        int       `json:"totalChats"`
	AvgResponseTimeMS float64   `json:"avgResponseTime"`
	SuccessRate       float64   `json:"successRate"`
	Conversions       int       `json:"conversions"`
}

// EventsResult wraps paginated events.
type EventsResult struct {
	Items []models.AnalyticsEvent `json:"items"`
	Total int64                   `json:"total"`
}

// Service coordinates analytics storage and aggregation.
type Service struct {
	repo          Repository
	retentionDays int
	enabled       bool
}

// NewService constructs a Service.
func NewService(repo Repository, retentionDays int, enabled bool) *Service {
	if retentionDays <= 0 {
		retentionDays = 90
	}
	return &Service{repo: repo, retentionDays: retentionDays, enabled: enabled}
}

// RecordChat persists chat metrics as analytics events.
func (s *Service) RecordChat(ctx context.Context, input RecordChatInput) error {
	if !s.enabled {
		return ErrAnalyticsDisabled
	}
	if input.ChatID == uuid.Nil {
		input.ChatID = uuid.New()
	}
	metadata := models.JSONB{"chat_id": input.ChatID.String()}
	if input.Metadata != nil {
		for k, v := range input.Metadata {
			metadata[k] = v
		}
	}

	event := models.AnalyticsEvent{
		Timestamp: input.Timestamp,
		EventType: "chat",
		Source:    emptyOrDefault(input.Source, "web"),
		Provider:  emptyOrDefault(input.Provider, "unknown"),
		Duration:  int(input.Duration.Milliseconds()),
		Success:   input.Success,
		UserAgent: input.UserAgent,
		Metadata:  metadata,
	}
	if _, err := s.repo.InsertEvent(ctx, event); err != nil {
		return err
	}
	return s.repo.UpsertSummary(ctx, event.Timestamp)
}

// RecordEvent persists a custom analytics event.
func (s *Service) RecordEvent(ctx context.Context, input RecordEventInput) error {
	if !s.enabled {
		return ErrAnalyticsDisabled
	}
	payload := input.Metadata
	if payload == nil {
		payload = models.JSONB{}
	}

	event := models.AnalyticsEvent{
		Timestamp: input.Timestamp,
		EventType: emptyOrDefault(input.Type, "custom"),
		Source:    emptyOrDefault(input.Source, "web"),
		Provider:  emptyOrDefault(input.Provider, "unknown"),
		Duration:  int(input.Duration.Milliseconds()),
		Success:   input.Success,
		UserAgent: input.UserAgent,
		Metadata:  payload,
	}
	if _, err := s.repo.InsertEvent(ctx, event); err != nil {
		return err
	}
	return s.repo.UpsertSummary(ctx, event.Timestamp)
}

// Summary fetches aggregated metrics for the requested period.
func (s *Service) Summary(ctx context.Context, filter RangeFilter) (SummaryRange, error) {
	if filter.Start.IsZero() {
		filter.Start = time.Now().AddDate(0, 0, -7)
	}
	if filter.End.IsZero() {
		filter.End = time.Now()
	}

	summary, err := s.repo.AggregateRange(ctx, filter)
	if err != nil {
		return SummaryRange{}, err
	}
	providers, err := s.repo.AggregateProviders(ctx, filter)
	if err != nil {
		return SummaryRange{}, err
	}
	daily, err := s.repo.AggregateDaily(ctx, filter)
	if err != nil {
		return SummaryRange{}, err
	}

	breakdown := make(map[string]ProviderSnapshot, len(providers))
	for _, item := range providers {
		breakdown[item.Provider] = ProviderSnapshot{
			TotalChats:        item.TotalChats,
			AvgResponseTimeMS: item.AvgResponseTimeMS,
			SuccessRate:       item.SuccessRate,
		}
	}

	daySeries := make([]DailySnapshot, 0, len(daily))
	for _, item := range daily {
		daySeries = append(daySeries, DailySnapshot{
			Date:              item.Day,
			TotalChats:        item.TotalChats,
			AvgResponseTimeMS: item.AvgResponseTimeMS,
			SuccessRate:       item.SuccessRate,
			Conversions:       item.Conversions,
		})
	}

	return SummaryRange{
		RangeStart:        filter.Start,
		RangeEnd:          filter.End,
		TotalChats:        summary.TotalChats,
		AvgResponseTimeMS: summary.AvgResponseTimeMS,
		SuccessRate:       summary.SuccessRate,
		UniqueUsers:       summary.UniqueUsers,
		Conversions:       summary.Conversions,
		ProviderBreakdown: breakdown,
		Daily:             daySeries,
	}, nil
}

// Events returns paginated analytics events for observability.
func (s *Service) Events(ctx context.Context, filter EventFilter) (EventsResult, error) {
	if filter.Limit <= 0 || filter.Limit > 500 {
		filter.Limit = 100
	}
	events, total, err := s.repo.ListEvents(ctx, filter)
	if err != nil {
		return EventsResult{}, err
	}
	return EventsResult{Items: events, Total: total}, nil
}

func emptyOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
