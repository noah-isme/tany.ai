package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

type stubRepository struct {
	inserted        []models.AnalyticsEvent
	summaries       []time.Time
	summary         SummaryAggregate
	providers       []ProviderAggregate
	daily           []DailyAggregate
	events          []models.AnalyticsEvent
	total           int64
	lastEventFilter EventFilter
}

func (s *stubRepository) InsertEvent(ctx context.Context, event models.AnalyticsEvent) (models.AnalyticsEvent, error) {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	s.inserted = append(s.inserted, event)
	return event, nil
}

func (s *stubRepository) ListEvents(ctx context.Context, filter EventFilter) ([]models.AnalyticsEvent, int64, error) {
	s.lastEventFilter = filter
	return s.events, s.total, nil
}

func (s *stubRepository) AggregateRange(ctx context.Context, filter RangeFilter) (SummaryAggregate, error) {
	return s.summary, nil
}

func (s *stubRepository) AggregateProviders(ctx context.Context, filter RangeFilter) ([]ProviderAggregate, error) {
	return s.providers, nil
}

func (s *stubRepository) AggregateDaily(ctx context.Context, filter RangeFilter) ([]DailyAggregate, error) {
	return s.daily, nil
}

func (s *stubRepository) UpsertSummary(ctx context.Context, date time.Time) error {
	s.summaries = append(s.summaries, date)
	return nil
}

func TestRecordChatDisabled(t *testing.T) {
	repo := &stubRepository{}
	service := NewService(repo, 30, false)
	err := service.RecordChat(context.Background(), RecordChatInput{})
	if err == nil || err != ErrAnalyticsDisabled {
		t.Fatalf("expected ErrAnalyticsDisabled, got %v", err)
	}
	if len(repo.inserted) != 0 {
		t.Fatalf("expected no events stored when disabled")
	}
}

func TestRecordChatInsertsEventAndSummary(t *testing.T) {
	repo := &stubRepository{}
	service := NewService(repo, 30, true)
	chatID := uuid.New()
	err := service.RecordChat(context.Background(), RecordChatInput{
		Timestamp: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
		Source:    "web",
		Provider:  "gemini",
		Duration:  250 * time.Millisecond,
		Success:   true,
		UserAgent: "test",
		ChatID:    chatID,
		Metadata: models.JSONB{
			"foo": "bar",
		},
	})
	if err != nil {
		t.Fatalf("record chat: %v", err)
	}
	if len(repo.inserted) != 1 {
		t.Fatalf("expected one event inserted, got %d", len(repo.inserted))
	}
	inserted := repo.inserted[0]
	if inserted.Provider != "gemini" {
		t.Fatalf("expected provider persisted, got %s", inserted.Provider)
	}
	if inserted.Duration != 250 {
		t.Fatalf("expected duration stored as milliseconds, got %d", inserted.Duration)
	}
	if len(repo.summaries) != 1 {
		t.Fatalf("expected summary upserted once, got %d", len(repo.summaries))
	}
}

func TestSummaryUsesRepositoryAggregates(t *testing.T) {
	repo := &stubRepository{
		summary: SummaryAggregate{
			TotalChats:        10,
			AvgResponseTimeMS: 120,
			SuccessRate:       0.9,
			UniqueUsers:       4,
			Conversions:       3,
		},
		providers: []ProviderAggregate{{
			Provider:          "gemini",
			TotalChats:        7,
			AvgResponseTimeMS: 100,
			SuccessRate:       0.95,
		}},
		daily: []DailyAggregate{{
			Day:               time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			TotalChats:        5,
			AvgResponseTimeMS: 110,
			SuccessRate:       0.8,
			Conversions:       2,
		}},
	}
	service := NewService(repo, 30, true)
	summary, err := service.Summary(context.Background(), RangeFilter{})
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.TotalChats != 10 {
		t.Fatalf("expected total chats to propagate")
	}
	if len(summary.ProviderBreakdown) != 1 {
		t.Fatalf("expected provider breakdown to contain entries")
	}
	if len(summary.Daily) != 1 {
		t.Fatalf("expected daily data to propagate")
	}
}

func TestEventsEnforcesLimit(t *testing.T) {
	repo := &stubRepository{
		events: []models.AnalyticsEvent{},
	}
	service := NewService(repo, 30, true)
	_, err := service.Events(context.Background(), EventFilter{Limit: 2000})
	if err != nil {
		t.Fatalf("events: %v", err)
	}
	if repo.lastEventFilter.Limit != 100 {
		t.Fatalf("expected limit to be capped to 100, got %d", repo.lastEventFilter.Limit)
	}
}
