package analytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// EventFilter captures optional filters for listing analytics events.
type EventFilter struct {
	Start    time.Time
	End      time.Time
	Source   string
	Provider string
	Limit    int
	Offset   int
	Type     string
}

// RangeFilter narrows down aggregation windows.
type RangeFilter struct {
	Start    time.Time
	End      time.Time
	Source   string
	Provider string
}

// ProviderAggregate summarises KPI per provider.
type ProviderAggregate struct {
	Provider          string  `db:"provider"`
	TotalChats        int     `db:"total_chats"`
	AvgResponseTimeMS float64 `db:"avg_response_time"`
	SuccessRate       float64 `db:"success_rate"`
}

// DailyAggregate summarises metrics per day.
type DailyAggregate struct {
	Day               time.Time `db:"bucket"`
	TotalChats        int       `db:"total_chats"`
	AvgResponseTimeMS float64   `db:"avg_response_time"`
	SuccessRate       float64   `db:"success_rate"`
	Conversions       int       `db:"conversions"`
}

// SummaryAggregate summarises totals for a period.
type SummaryAggregate struct {
	TotalChats        int     `db:"total_chats"`
	AvgResponseTimeMS float64 `db:"avg_response_time"`
	SuccessRate       float64 `db:"success_rate"`
	UniqueUsers       int     `db:"unique_users"`
	Conversions       int     `db:"conversions"`
}

// Repository persists analytics related data.
type Repository interface {
	InsertEvent(ctx context.Context, event models.AnalyticsEvent) (models.AnalyticsEvent, error)
	ListEvents(ctx context.Context, filter EventFilter) ([]models.AnalyticsEvent, int64, error)
	AggregateRange(ctx context.Context, filter RangeFilter) (SummaryAggregate, error)
	AggregateProviders(ctx context.Context, filter RangeFilter) ([]ProviderAggregate, error)
	AggregateDaily(ctx context.Context, filter RangeFilter) ([]DailyAggregate, error)
	UpsertSummary(ctx context.Context, date time.Time) error
}

// NewRepository constructs a SQL backed analytics repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) InsertEvent(ctx context.Context, event models.AnalyticsEvent) (models.AnalyticsEvent, error) {
	const query = `INSERT INTO analytics_events (timestamp, event_type, source, provider, duration_ms, success, user_agent, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, timestamp, event_type, source, provider, duration_ms, success, user_agent, metadata`

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Metadata == nil {
		event.Metadata = models.JSONB{}
	}

	var created models.AnalyticsEvent
	if err := r.db.GetContext(ctx, &created, query,
		event.Timestamp,
		event.EventType,
		event.Source,
		event.Provider,
		event.Duration,
		event.Success,
		event.UserAgent,
		event.Metadata,
	); err != nil {
		return models.AnalyticsEvent{}, err
	}
	return created, nil
}

func (r *repository) ListEvents(ctx context.Context, filter EventFilter) ([]models.AnalyticsEvent, int64, error) {
	base := strings.Builder{}
	base.WriteString("FROM analytics_events WHERE 1=1")
	args := make([]interface{}, 0, 6)

	addFilter := func(clause string, value interface{}) {
		args = append(args, value)
		base.WriteString(" AND ")
		base.WriteString(fmt.Sprintf(clause, len(args)))
	}

	if !filter.Start.IsZero() {
		addFilter("timestamp >= $%d", filter.Start)
	}
	if !filter.End.IsZero() {
		addFilter("timestamp <= $%d", filter.End)
	}
	if filter.Source != "" {
		addFilter("source = $%d", filter.Source)
	}
	if filter.Provider != "" {
		addFilter("provider = $%d", filter.Provider)
	}
	if filter.Type != "" {
		addFilter("event_type = $%d", filter.Type)
	}

	query := fmt.Sprintf("SELECT id, timestamp, event_type, source, provider, duration_ms, success, user_agent, metadata %s ORDER BY timestamp DESC", base.String())

	limit := 100
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}
	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)

	rows := make([]models.AnalyticsEvent, 0, limit)
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", base.String())
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *repository) AggregateRange(ctx context.Context, filter RangeFilter) (SummaryAggregate, error) {
	const base = `SELECT
    COUNT(*) FILTER (WHERE event_type = 'chat') AS total_chats,
    COALESCE(AVG(duration_ms) FILTER (WHERE event_type = 'chat'), 0) AS avg_response_time,
    COALESCE(AVG(CASE WHEN success THEN 1 ELSE 0 END) FILTER (WHERE event_type = 'chat'), 0) AS success_rate,
    COALESCE(COUNT(DISTINCT metadata->>'chat_id') FILTER (WHERE event_type = 'chat'), 0) AS unique_users,
    COUNT(*) FILTER (WHERE event_type = 'lead') AS conversions
FROM analytics_events
WHERE 1=1`

	query := strings.Builder{}
	query.WriteString(base)
	args := make([]interface{}, 0, 4)
	add := func(clause string, value interface{}) {
		args = append(args, value)
		query.WriteString(" AND ")
		query.WriteString(fmt.Sprintf(clause, len(args)))
	}

	if !filter.Start.IsZero() {
		add("timestamp >= $%d", filter.Start)
	}
	if !filter.End.IsZero() {
		add("timestamp <= $%d", filter.End)
	}
	if filter.Source != "" {
		add("source = $%d", filter.Source)
	}
	if filter.Provider != "" {
		add("provider = $%d", filter.Provider)
	}

	var agg SummaryAggregate
	if err := r.db.GetContext(ctx, &agg, query.String(), args...); err != nil {
		return SummaryAggregate{}, err
	}
	return agg, nil
}

func (r *repository) AggregateProviders(ctx context.Context, filter RangeFilter) ([]ProviderAggregate, error) {
	const base = `SELECT
    provider,
    COUNT(*) FILTER (WHERE event_type = 'chat') AS total_chats,
    COALESCE(AVG(duration_ms) FILTER (WHERE event_type = 'chat'), 0) AS avg_response_time,
    COALESCE(AVG(CASE WHEN success THEN 1 ELSE 0 END) FILTER (WHERE event_type = 'chat'), 0) AS success_rate
FROM analytics_events
WHERE event_type = 'chat'`

	query := strings.Builder{}
	query.WriteString(base)
	args := make([]interface{}, 0, 4)
	add := func(clause string, value interface{}) {
		args = append(args, value)
		query.WriteString(" AND ")
		query.WriteString(fmt.Sprintf(clause, len(args)))
	}
	if !filter.Start.IsZero() {
		add("timestamp >= $%d", filter.Start)
	}
	if !filter.End.IsZero() {
		add("timestamp <= $%d", filter.End)
	}
	if filter.Source != "" {
		add("source = $%d", filter.Source)
	}
	if filter.Provider != "" {
		add("provider = $%d", filter.Provider)
	}
	query.WriteString(" GROUP BY provider ORDER BY provider")

	rows := []ProviderAggregate{}
	if err := r.db.SelectContext(ctx, &rows, query.String(), args...); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) AggregateDaily(ctx context.Context, filter RangeFilter) ([]DailyAggregate, error) {
	const base = `SELECT
    date_trunc('day', timestamp) AS bucket,
    COUNT(*) FILTER (WHERE event_type = 'chat') AS total_chats,
    COALESCE(AVG(duration_ms) FILTER (WHERE event_type = 'chat'), 0) AS avg_response_time,
    COALESCE(AVG(CASE WHEN success THEN 1 ELSE 0 END) FILTER (WHERE event_type = 'chat'), 0) AS success_rate,
    COUNT(*) FILTER (WHERE event_type = 'lead') AS conversions
FROM analytics_events
WHERE 1=1`

	query := strings.Builder{}
	query.WriteString(base)
	args := make([]interface{}, 0, 4)
	add := func(clause string, value interface{}) {
		args = append(args, value)
		query.WriteString(" AND ")
		query.WriteString(fmt.Sprintf(clause, len(args)))
	}
	if !filter.Start.IsZero() {
		add("timestamp >= $%d", filter.Start)
	}
	if !filter.End.IsZero() {
		add("timestamp <= $%d", filter.End)
	}
	if filter.Source != "" {
		add("source = $%d", filter.Source)
	}
	if filter.Provider != "" {
		add("provider = $%d", filter.Provider)
	}
	query.WriteString(" GROUP BY bucket ORDER BY bucket")

	rows := []DailyAggregate{}
	if err := r.db.SelectContext(ctx, &rows, query.String(), args...); err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) UpsertSummary(ctx context.Context, date time.Time) error {
	const query = `WITH provider_stats AS (
    SELECT
        provider,
        COUNT(*) FILTER (WHERE event_type = 'chat') AS total_chats,
        COALESCE(AVG(duration_ms) FILTER (WHERE event_type = 'chat'), 0) AS avg_response_time,
        COALESCE(AVG(CASE WHEN success THEN 1 ELSE 0 END) FILTER (WHERE event_type = 'chat'), 0) AS success_rate
    FROM analytics_events
    WHERE DATE_TRUNC('day', timestamp)::date = $1::date
      AND provider IS NOT NULL
    GROUP BY provider
), agg AS (
    SELECT
        DATE_TRUNC('day', timestamp)::date AS day,
        COUNT(*) FILTER (WHERE event_type = 'chat') AS total_chats,
        COALESCE(AVG(duration_ms) FILTER (WHERE event_type = 'chat'), 0) AS avg_response_time,
        COALESCE(AVG(CASE WHEN success THEN 1 ELSE 0 END) FILTER (WHERE event_type = 'chat'), 0) AS success_rate,
        COALESCE(COUNT(DISTINCT metadata->>'chat_id') FILTER (WHERE event_type = 'chat'), 0) AS unique_users,
        COUNT(*) FILTER (WHERE event_type = 'lead') AS conversions,
        COALESCE(
            (
                SELECT jsonb_object_agg(provider, jsonb_build_object(
                    'totalChats', total_chats,
                    'avgResponseTime', avg_response_time,
                    'successRate', success_rate
                )) FROM provider_stats
            ), '{}'::jsonb
        ) AS provider_breakdown
    FROM analytics_events
    WHERE DATE_TRUNC('day', timestamp)::date = $1::date
    GROUP BY day
)
INSERT INTO analytics_summary (date, total_chats, avg_response_time, success_rate, unique_users, conversions, provider_breakdown)
SELECT day, total_chats, avg_response_time, success_rate, unique_users, conversions, provider_breakdown
FROM agg
ON CONFLICT (date)
DO UPDATE SET
    total_chats = EXCLUDED.total_chats,
    avg_response_time = EXCLUDED.avg_response_time,
    success_rate = EXCLUDED.success_rate,
    unique_users = EXCLUDED.unique_users,
    conversions = EXCLUDED.conversions,
    provider_breakdown = EXCLUDED.provider_breakdown,
    updated_at = NOW();`

	if date.IsZero() {
		date = time.Now()
	}
	if _, err := r.db.ExecContext(ctx, query, date); err != nil {
		return err
	}
	return nil
}
