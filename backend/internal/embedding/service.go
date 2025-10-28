package embedding

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

// Provider defines the behaviour required from an embedding generator.
type Provider interface {
	Embed(ctx context.Context, input string) ([]float32, error)
}

// ErrProviderMissing indicates an embedding provider was not configured.
var ErrProviderMissing = errors.New("embedding provider is not configured")

// Summary captures high level metrics for admin dashboards.
type Summary struct {
	Enabled         bool       `json:"enabled"`
	Provider        string     `json:"provider"`
	Dimension       int        `json:"dimension"`
	Count           int64      `json:"count"`
	Weight          float64    `json:"weight"`
	LastReindexedAt *time.Time `json:"lastReindexedAt,omitempty"`
	LastResetAt     *time.Time `json:"lastResetAt,omitempty"`
}

// Snippet represents a personalization context fragment returned to the prompt builder.
type Snippet struct {
	Kind     string         `json:"kind"`
	Score    float64        `json:"score"`
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata"`
}

// PersonalizationResult returns personalization context for chat prompts.
type PersonalizationResult struct {
	Enabled  bool      `json:"enabled"`
	Weight   float64   `json:"weight"`
	Provider string    `json:"provider"`
	Snippets []Snippet `json:"snippets"`
}

// Options configures the service behaviour.
type Options struct {
	Enabled       bool
	ProviderName  string
	Dimension     int
	CacheTTL      time.Duration
	DefaultWeight float64
	MinScore      float64
	MaxSnippets   int
}

// Service orchestrates embedding management and personalization retrieval.
type Service struct {
	repo     Repository
	provider Provider

	enabled      bool
	providerName string
	dimension    int

	cacheTTL    time.Duration
	minScore    float64
	maxSnippets int

	mu     sync.RWMutex
	cache  map[string]cacheEntry
	config models.EmbeddingConfig
	weight float64
}

type cacheEntry struct {
	snippets []Snippet
	expires  time.Time
}

// NewService constructs a Service using the provided repository and provider.
func NewService(ctx context.Context, repo Repository, provider Provider, opts Options) (*Service, error) {
	if repo == nil {
		return nil, errors.New("repository is required")
	}
	if opts.Dimension <= 0 {
		opts.Dimension = 1536
	}
	if opts.CacheTTL <= 0 {
		opts.CacheTTL = 24 * time.Hour
	}
	if opts.DefaultWeight <= 0 {
		opts.DefaultWeight = 0.65
	}
	if opts.MinScore <= 0 {
		opts.MinScore = 0.45
	}
	if opts.MaxSnippets <= 0 {
		opts.MaxSnippets = 5
	}

	cfg, err := repo.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}
	weight := clampWeight(cfg.Weight)
	if weight == 0 {
		weight = clampWeight(opts.DefaultWeight)
	}

	service := &Service{
		repo:         repo,
		provider:     provider,
		enabled:      opts.Enabled,
		providerName: strings.TrimSpace(opts.ProviderName),
		dimension:    opts.Dimension,
		cacheTTL:     opts.CacheTTL,
		minScore:     opts.MinScore,
		maxSnippets:  opts.MaxSnippets,
		cache:        make(map[string]cacheEntry),
		config:       cfg,
		weight:       weight,
	}
	return service, nil
}

// Enabled reports whether personalization is globally enabled.
func (s *Service) Enabled() bool {
	return s != nil && s.enabled && s.provider != nil
}

// Weight returns the personalization weight multiplier.
func (s *Service) Weight() float64 {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.weight
}

// Summary returns high level metrics for the admin UI.
func (s *Service) Summary(ctx context.Context) (Summary, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return Summary{}, err
	}
	s.mu.RLock()
	cfg := s.config
	weight := s.weight
	s.mu.RUnlock()

	summary := Summary{
		Enabled:   s.Enabled(),
		Provider:  s.providerName,
		Dimension: s.dimension,
		Count:     count,
		Weight:    weight,
	}
	if !cfg.LastReindexedAt.IsZero() {
		ts := cfg.LastReindexedAt.UTC()
		summary.LastReindexedAt = &ts
	}
	if !cfg.LastResetAt.IsZero() {
		ts := cfg.LastResetAt.UTC()
		summary.LastResetAt = &ts
	}
	return summary, nil
}

// SetWeight updates personalization weight and persists it.
func (s *Service) SetWeight(ctx context.Context, weight float64) (float64, error) {
	if s == nil {
		return 0, errors.New("service is nil")
	}
	clamped := clampWeight(weight)
	s.mu.Lock()
	s.weight = clamped
	s.config.Weight = clamped
	cfg := s.config
	s.mu.Unlock()

	if err := s.repo.SaveConfig(ctx, cfg); err != nil {
		return 0, err
	}
	s.invalidateCache()
	return clamped, nil
}

// Reset removes all embeddings and records the reset timestamp.
func (s *Service) Reset(ctx context.Context) error {
	if err := s.repo.DeleteAll(ctx); err != nil {
		return err
	}
	now := time.Now().UTC()
	s.mu.Lock()
	s.config.LastResetAt = now
	s.config.LastReindexedAt = time.Time{}
	cfg := s.config
	s.mu.Unlock()

	s.invalidateCache()
	return s.repo.SaveConfig(ctx, cfg)
}

// Reindex rebuilds embeddings using the provided knowledge base snapshot.
func (s *Service) Reindex(ctx context.Context, base kb.KnowledgeBase) (int, error) {
	if s.provider == nil {
		return 0, ErrProviderMissing
	}
	entries := gatherEntries(base)
	if len(entries) == 0 {
		return 0, errors.New("no content available for embeddings")
	}
	if err := s.repo.DeleteAll(ctx); err != nil {
		return 0, err
	}

	created := 0
	for _, entry := range entries {
		vector, err := s.provider.Embed(ctx, entry.Content)
		if err != nil {
			return created, fmt.Errorf("generate embedding for %s: %w", entry.Kind, err)
		}
		metadata := models.JSONB(entry.Metadata)
		if metadata == nil {
			metadata = models.JSONB{}
		}
		embedding := models.Embedding{
			ID:       uuid.New(),
			Kind:     entry.Kind,
			RefID:    entry.RefID,
			Content:  entry.Content,
			Vector:   vector,
			Metadata: metadata,
		}
		if err := s.repo.Upsert(ctx, embedding); err != nil {
			return created, err
		}
		created++
	}

	now := time.Now().UTC()
	s.mu.Lock()
	s.config.LastReindexedAt = now
	cfg := s.config
	s.mu.Unlock()

	s.invalidateCache()
	if err := s.repo.SaveConfig(ctx, cfg); err != nil {
		return created, err
	}
	return created, nil
}

// Personalize performs a similarity lookup to enrich prompts.
func (s *Service) Personalize(ctx context.Context, question string) (PersonalizationResult, error) {
	result := PersonalizationResult{
		Enabled:  s.Enabled(),
		Provider: s.providerName,
	}
	if !result.Enabled {
		result.Weight = s.Weight()
		return result, nil
	}

	question = strings.TrimSpace(question)
	if question == "" {
		result.Weight = s.Weight()
		return result, nil
	}

	cacheKey := cacheKey(question)
	if snippets, ok := s.lookupCache(cacheKey); ok {
		result.Snippets = snippets
		result.Weight = s.Weight()
		return result, nil
	}

	vector, err := s.provider.Embed(ctx, question)
	if err != nil {
		return result, err
	}

	matches, err := s.repo.Similar(ctx, vector, s.maxSnippets, s.minScore, nil)
	if err != nil {
		return result, err
	}
	snippets := make([]Snippet, 0, len(matches))
	for _, match := range matches {
		snippets = append(snippets, Snippet{
			Kind:    match.Kind,
			Score:   match.Score,
			Content: match.Content,
			Metadata: func(data models.JSONB) map[string]any {
				if data == nil {
					return map[string]any{}
				}
				clone := make(map[string]any, len(data))
				for k, v := range data {
					clone[k] = v
				}
				return clone
			}(match.Metadata),
		})
	}
	sort.SliceStable(snippets, func(i, j int) bool {
		if snippets[i].Score == snippets[j].Score {
			return snippets[i].Kind < snippets[j].Kind
		}
		return snippets[i].Score > snippets[j].Score
	})

	s.storeCache(cacheKey, snippets)
	result.Snippets = snippets
	result.Weight = s.Weight()
	return result, nil
}

func (s *Service) lookupCache(key string) ([]Snippet, bool) {
	s.mu.RLock()
	entry, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		return nil, false
	}
	return entry.snippets, true
}

func (s *Service) storeCache(key string, snippets []Snippet) {
	if len(snippets) == 0 {
		return
	}
	s.mu.Lock()
	s.cache[key] = cacheEntry{snippets: snippets, expires: time.Now().Add(s.cacheTTL)}
	s.mu.Unlock()
}

func (s *Service) invalidateCache() {
	s.mu.Lock()
	s.cache = make(map[string]cacheEntry)
	s.mu.Unlock()
}

func cacheKey(question string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(question)))
	return hex.EncodeToString(sum[:])
}

func clampWeight(weight float64) float64 {
	if weight < 0 {
		return 0
	}
	if weight > 1 {
		return 1
	}
	return weight
}

type entry struct {
	Kind     string
	RefID    *uuid.UUID
	Content  string
	Metadata map[string]any
}

func gatherEntries(base kb.KnowledgeBase) []entry {
	entries := make([]entry, 0, 16)
	profileContent := buildProfileContent(base.Profile)
	if profileContent != "" {
		entries = append(entries, entry{
			Kind:    "profile",
			Content: profileContent,
			Metadata: map[string]any{
				"name":     base.Profile.Name,
				"title":    base.Profile.Title,
				"location": base.Profile.Location,
			},
		})
	}
	for _, service := range base.Services {
		if strings.TrimSpace(service.Name) == "" {
			continue
		}
		var refID *uuid.UUID
		if id, err := uuid.Parse(service.ID); err == nil {
			refID = &id
		}
		entries = append(entries, entry{
			Kind:    "service",
			RefID:   refID,
			Content: buildServiceContent(service),
			Metadata: map[string]any{
				"name":        service.Name,
				"description": service.Description,
				"order":       service.Order,
			},
		})
	}
	for _, project := range base.Projects {
		if strings.TrimSpace(project.Title) == "" {
			continue
		}
		var refID *uuid.UUID
		if id, err := uuid.Parse(project.ID); err == nil {
			refID = &id
		}
		entries = append(entries, entry{
			Kind:    "project",
			RefID:   refID,
			Content: buildProjectContent(project),
			Metadata: map[string]any{
				"title":    project.Title,
				"category": project.Category,
			},
		})
	}
	for _, post := range base.Posts {
		if strings.TrimSpace(post.Title) == "" {
			continue
		}
		var refID *uuid.UUID
		if id, err := uuid.Parse(post.ID); err == nil {
			refID = &id
		}
		entries = append(entries, entry{
			Kind:    "post",
			RefID:   refID,
			Content: buildPostContent(post),
			Metadata: map[string]any{
				"title":  post.Title,
				"source": post.Source,
			},
		})
	}
	return entries
}

func buildProfileContent(profile kb.Profile) string {
	var parts []string
	if strings.TrimSpace(profile.Name) != "" {
		parts = append(parts, fmt.Sprintf("Nama: %s", profile.Name))
	}
	if strings.TrimSpace(profile.Title) != "" {
		parts = append(parts, fmt.Sprintf("Peran: %s", profile.Title))
	}
	if strings.TrimSpace(profile.Location) != "" {
		parts = append(parts, fmt.Sprintf("Berbasis di %s", profile.Location))
	}
	if strings.TrimSpace(profile.Bio) != "" {
		parts = append(parts, fmt.Sprintf("Bio: %s", profile.Bio))
	}
	if len(parts) == 0 {
		return ""
	}
	return "Persona profesional. " + strings.Join(parts, ". ") + "."
}

func buildServiceContent(service kb.Service) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Layanan: %s", strings.TrimSpace(service.Name)))
	if strings.TrimSpace(service.Description) != "" {
		parts = append(parts, fmt.Sprintf("Deskripsi: %s", service.Description))
	}
	if len(service.PriceRange) > 0 {
		parts = append(parts, fmt.Sprintf("Rentang harga: %s", strings.Join(service.PriceRange, " - ")))
	}
	if strings.TrimSpace(service.DurationLabel) != "" {
		parts = append(parts, fmt.Sprintf("Durasi: %s", service.DurationLabel))
	}
	return strings.Join(parts, ". ") + ". Fokuskan pada manfaat dan hasil nyata."
}

func buildProjectContent(project kb.Project) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Proyek: %s", strings.TrimSpace(project.Title)))
	if strings.TrimSpace(project.Description) != "" {
		parts = append(parts, fmt.Sprintf("Ringkasan: %s", project.Description))
	}
	if len(project.TechStack) > 0 {
		parts = append(parts, fmt.Sprintf("Teknologi: %s", strings.Join(project.TechStack, ", ")))
	}
	if strings.TrimSpace(project.ProjectURL) != "" {
		parts = append(parts, fmt.Sprintf("URL: %s", project.ProjectURL))
	}
	if project.IsFeatured {
		parts = append(parts, "Prioritas tinggi, tampilkan dengan kebanggaan dan kepercayaan diri")
	}
	return strings.Join(parts, ". ") + ". Tekankan storytelling yang hangat dan solutif."
}

func buildPostContent(post kb.Post) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Tulisan: %s", strings.TrimSpace(post.Title)))
	if strings.TrimSpace(post.Summary) != "" {
		parts = append(parts, fmt.Sprintf("Inti: %s", post.Summary))
	}
	if strings.TrimSpace(post.Source) != "" {
		parts = append(parts, fmt.Sprintf("Sumber: %s", post.Source))
	}
	if !post.PublishedAt.IsZero() {
		parts = append(parts, fmt.Sprintf("Terbit: %s", post.PublishedAt.Format("2006-01-02")))
	}
	return strings.Join(parts, ". ") + ". Ambil tone yang reflektif dan humanis."
}
