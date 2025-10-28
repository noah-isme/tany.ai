package kb

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type cacheEntry struct {
	data    KnowledgeBase
	etag    string
	expires time.Time
}

// Aggregator loads and caches knowledge base data from the database.
type Aggregator struct {
	db  *sqlx.DB
	ttl time.Duration

	mu    sync.RWMutex
	cache *cacheEntry
}

// NewAggregator constructs a new Aggregator with the provided cache TTL.
func NewAggregator(db *sqlx.DB, ttl time.Duration) *Aggregator {
	if ttl <= 0 {
		ttl = time.Minute
	}
	return &Aggregator{db: db, ttl: ttl}
}

// Get retrieves the knowledge base, optionally serving it from cache.
// It returns the data, the computed ETag and a boolean indicating a cache hit.
func (a *Aggregator) Get(ctx context.Context) (KnowledgeBase, string, bool, error) {
	now := time.Now()

	a.mu.RLock()
	entry := a.cache
	if entry != nil && now.Before(entry.expires) {
		data := entry.data
		etag := entry.etag
		a.mu.RUnlock()
		return data, etag, true, nil
	}
	a.mu.RUnlock()

	data, err := a.load(ctx)
	if err != nil {
		return KnowledgeBase{}, "", false, err
	}
	etag, err := computeETag(data)
	if err != nil {
		return KnowledgeBase{}, "", false, err
	}

	a.mu.Lock()
	a.cache = &cacheEntry{data: data, etag: etag, expires: now.Add(a.ttl)}
	a.mu.Unlock()

	return data, etag, false, nil
}

// Invalidate clears the in-memory cache so subsequent Get calls refetch data.
func (a *Aggregator) Invalidate() {
	a.mu.Lock()
	a.cache = nil
	a.mu.Unlock()
}

// CacheTTL returns the configured cache duration.
func (a *Aggregator) CacheTTL() time.Duration {
	return a.ttl
}

func (a *Aggregator) load(ctx context.Context) (KnowledgeBase, error) {
	profile, err := a.fetchProfile(ctx)
	if err != nil {
		return KnowledgeBase{}, err
	}

	skills, err := a.fetchSkills(ctx)
	if err != nil {
		return KnowledgeBase{}, err
	}

	services, err := a.fetchServices(ctx)
	if err != nil {
		return KnowledgeBase{}, err
	}

	projects, err := a.fetchProjects(ctx)
	if err != nil {
		return KnowledgeBase{}, err
	}

	return KnowledgeBase{Profile: profile, Skills: skills, Services: services, Projects: projects}, nil
}

func (a *Aggregator) fetchProfile(ctx context.Context) (Profile, error) {
	const query = `SELECT id, name, title, bio, email, phone, location, avatar_url, updated_at FROM profile ORDER BY updated_at DESC LIMIT 1`

	var row struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		Title     *string   `db:"title"`
		Bio       *string   `db:"bio"`
		Email     *string   `db:"email"`
		Phone     *string   `db:"phone"`
		Location  *string   `db:"location"`
		AvatarURL *string   `db:"avatar_url"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := a.db.GetContext(ctx, &row, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, nil
		}
		return Profile{}, err
	}

	return Profile{
		Name:      row.Name,
		Title:     deref(row.Title),
		Bio:       deref(row.Bio),
		Email:     deref(row.Email),
		Phone:     deref(row.Phone),
		Location:  deref(row.Location),
		AvatarURL: deref(row.AvatarURL),
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (a *Aggregator) fetchSkills(ctx context.Context) ([]Skill, error) {
	const query = `SELECT name FROM skills ORDER BY "order" ASC, name ASC`
	var rows []struct {
		Name string `db:"name"`
	}
	if err := a.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	skills := make([]Skill, 0, len(rows))
	for _, row := range rows {
		skills = append(skills, Skill{Name: row.Name})
	}
	return skills, nil
}

func (a *Aggregator) fetchServices(ctx context.Context) ([]Service, error) {
	const query = `SELECT id, name, description, price_min, price_max, currency, duration_label, "order" FROM services WHERE is_active = TRUE ORDER BY "order" ASC, name ASC`
	var rows []struct {
		ID            uuid.UUID `db:"id"`
		Name          string    `db:"name"`
		Description   *string   `db:"description"`
		PriceMin      *float64  `db:"price_min"`
		PriceMax      *float64  `db:"price_max"`
		Currency      *string   `db:"currency"`
		DurationLabel *string   `db:"duration_label"`
		Order         int       `db:"order"`
	}
	if err := a.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	services := make([]Service, 0, len(rows))
	for _, row := range rows {
		services = append(services, Service{
			ID:            row.ID.String(),
			Name:          row.Name,
			Description:   deref(row.Description),
			Currency:      deref(row.Currency),
			DurationLabel: deref(row.DurationLabel),
			PriceRange:    formatPriceRange(row.PriceMin, row.PriceMax, deref(row.Currency)),
			Order:         row.Order,
		})
	}
	return services, nil
}

func (a *Aggregator) fetchProjects(ctx context.Context) ([]Project, error) {
	const query = `SELECT id, title, description, tech_stack, project_url, category, duration_label, price_label, budget_label, "order", is_featured FROM projects ORDER BY is_featured DESC, "order" ASC, title ASC`
	var rows []struct {
		ID            uuid.UUID      `db:"id"`
		Title         string         `db:"title"`
		Description   *string        `db:"description"`
		TechStack     pq.StringArray `db:"tech_stack"`
		ProjectURL    *string        `db:"project_url"`
		Category      *string        `db:"category"`
		DurationLabel *string        `db:"duration_label"`
		PriceLabel    *string        `db:"price_label"`
		BudgetLabel   *string        `db:"budget_label"`
		Order         int            `db:"order"`
		IsFeatured    bool           `db:"is_featured"`
	}
	if err := a.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	projects := make([]Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, Project{
			ID:            row.ID.String(),
			Title:         row.Title,
			Description:   deref(row.Description),
			TechStack:     []string(row.TechStack),
			ProjectURL:    deref(row.ProjectURL),
			Category:      deref(row.Category),
			DurationLabel: deref(row.DurationLabel),
			PriceLabel:    deref(row.PriceLabel),
			BudgetLabel:   deref(row.BudgetLabel),
			IsFeatured:    row.IsFeatured,
			Order:         row.Order,
		})
	}
	return projects, nil
}

func computeETag(data KnowledgeBase) (string, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "W/\"" + hex.EncodeToString(sum[:]) + "\"", nil
}

func deref[T any](ptr *T) T {
	var zero T
	if ptr == nil {
		return zero
	}
	return *ptr
}

func formatPriceRange(min, max *float64, currency string) []string {
	if min == nil && max == nil {
		return nil
	}

	format := func(value *float64) string {
		if value == nil {
			return ""
		}
		return formatPrice(*value, currency)
	}

	minStr := format(min)
	maxStr := format(max)
	values := make([]string, 0, 2)
	if minStr != "" {
		values = append(values, minStr)
	}
	if maxStr != "" && maxStr != minStr {
		values = append(values, maxStr)
	}
	return values
}

func formatPrice(value float64, currency string) string {
	if currency == "" {
		return formatFloat(value)
	}
	return currency + " " + formatFloat(value)
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
