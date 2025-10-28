package embedding

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

type memoryRepo struct {
	embeddings []models.Embedding
	config     models.EmbeddingConfig
}

func (m *memoryRepo) Upsert(_ context.Context, embedding models.Embedding) error {
	for i, existing := range m.embeddings {
		if existing.ID == embedding.ID {
			m.embeddings[i] = embedding
			return nil
		}
	}
	m.embeddings = append(m.embeddings, embedding)
	return nil
}

func (m *memoryRepo) DeleteByKind(_ context.Context, kind string) error {
	filtered := m.embeddings[:0]
	for _, e := range m.embeddings {
		if e.Kind != kind {
			filtered = append(filtered, e)
		}
	}
	m.embeddings = filtered
	return nil
}

func (m *memoryRepo) DeleteAll(context.Context) error {
	m.embeddings = nil
	return nil
}

func (m *memoryRepo) Count(context.Context) (int64, error) {
	return int64(len(m.embeddings)), nil
}

func (m *memoryRepo) Similar(_ context.Context, _ []float32, limit int, minScore float64, kinds []string) ([]models.EmbeddingMatch, error) {
	matches := make([]models.EmbeddingMatch, 0, len(m.embeddings))
	for _, e := range m.embeddings {
		if len(kinds) > 0 && !contains(kinds, e.Kind) {
			continue
		}
		matches = append(matches, models.EmbeddingMatch{
			ID:       e.ID,
			Kind:     e.Kind,
			RefID:    e.RefID,
			Content:  e.Content,
			Metadata: e.Metadata,
			Score:    0.9,
		})
	}
	if minScore > 0 {
		filtered := matches[:0]
		for _, match := range matches {
			if match.Score >= minScore {
				filtered = append(filtered, match)
			}
		}
		matches = filtered
	}
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	return matches, nil
}

func (m *memoryRepo) LoadConfig(context.Context) (models.EmbeddingConfig, error) {
	if m.config.Weight == 0 {
		m.config.Weight = 0.65
	}
	return m.config, nil
}

func (m *memoryRepo) SaveConfig(_ context.Context, cfg models.EmbeddingConfig) error {
	m.config = cfg
	return nil
}

type stubProvider struct {
	calls int
	vec   []float32
}

func (s *stubProvider) Embed(_ context.Context, _ string) ([]float32, error) {
	s.calls++
	return append([]float32(nil), s.vec...), nil
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func TestServiceSetWeightPersists(t *testing.T) {
	repo := &memoryRepo{}
	provider := &stubProvider{vec: []float32{0, 0, 1}}
	ctx := context.Background()

	svc, err := NewService(ctx, repo, provider, Options{Enabled: true, Dimension: 3, DefaultWeight: 0.5})
	require.NoError(t, err)

	updated, err := svc.SetWeight(ctx, 0.82)
	require.NoError(t, err)
	require.InDelta(t, 0.82, updated, 1e-9)
	require.InDelta(t, 0.82, repo.config.Weight, 1e-9)
}

func TestServicePersonalizeCachesResults(t *testing.T) {
	repo := &memoryRepo{}
	provider := &stubProvider{vec: []float32{0.1, 0.2, 0.3}}
	ctx := context.Background()

	svc, err := NewService(ctx, repo, provider, Options{Enabled: true, Dimension: 3, CacheTTL: time.Hour})
	require.NoError(t, err)

	// Seed repository with one embedding so Similar returns data.
	embedding := models.Embedding{
		ID:      uuid.New(),
		Kind:    "profile",
		Content: "Persona profesional.",
		Metadata: models.JSONB{
			"name": "Noah",
		},
	}
	require.NoError(t, repo.Upsert(ctx, embedding))

	result, err := svc.Personalize(ctx, "Apa gaya komunikasi kamu?")
	require.NoError(t, err)
	require.True(t, result.Enabled)
	require.Len(t, result.Snippets, 1)
	require.Equal(t, 1, provider.calls)

	// Second call should hit cache.
	result, err = svc.Personalize(ctx, "Apa gaya komunikasi kamu?")
	require.NoError(t, err)
	require.Len(t, result.Snippets, 1)
	require.Equal(t, 1, provider.calls)
}

func TestServiceReindexCreatesEmbeddings(t *testing.T) {
	repo := &memoryRepo{}
	provider := &stubProvider{vec: []float32{0.1, 0.2, 0.3}}
	ctx := context.Background()

	svc, err := NewService(ctx, repo, provider, Options{Enabled: true, Dimension: 3})
	require.NoError(t, err)

	base := kb.KnowledgeBase{
		Profile:  kb.Profile{Name: "Noah", Title: "AI Engineer", Bio: "Suka menulis dengan tone hangat."},
		Services: []kb.Service{{ID: uuid.New().String(), Name: "AI Consultation", Description: "Diskusi strategi AI."}},
	}

	count, err := svc.Reindex(ctx, base)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	stored, err := repo.Count(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 2, stored)
	require.False(t, repo.config.LastReindexedAt.IsZero())
}
