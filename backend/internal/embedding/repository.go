package embedding

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

const personalizationConfigKey = "personalization"

// Repository defines persistence operations for embeddings and personalization settings.
type Repository interface {
	Upsert(ctx context.Context, embedding models.Embedding) error
	DeleteByKind(ctx context.Context, kind string) error
	DeleteAll(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Similar(ctx context.Context, vector []float32, limit int, minScore float64, kinds []string) ([]models.EmbeddingMatch, error)
	LoadConfig(ctx context.Context) (models.EmbeddingConfig, error)
	SaveConfig(ctx context.Context, cfg models.EmbeddingConfig) error
}

// NewRepository builds a SQL-backed Repository implementation.
func NewRepository(db *sqlx.DB, dimension int) Repository {
	if dimension <= 0 {
		dimension = 1536
	}
	return &sqlRepository{db: db, dimension: dimension}
}

type sqlRepository struct {
	db        *sqlx.DB
	dimension int
}

func (r *sqlRepository) Upsert(ctx context.Context, embedding models.Embedding) error {
	if embedding.ID == uuid.Nil {
		embedding.ID = uuid.New()
	}
	if embedding.Metadata == nil {
		embedding.Metadata = models.JSONB{}
	}
	if embedding.Content == "" {
		return errors.New("embedding content is required")
	}
	var vector any
	if len(embedding.Vector) > 0 {
		if len(embedding.Vector) != r.dimension {
			return fmt.Errorf("embedding dimension mismatch: expected %d got %d", r.dimension, len(embedding.Vector))
		}
		vector = pgvector.NewVector(embedding.Vector)
	}

	_, err := r.db.ExecContext(ctx, `
        INSERT INTO embeddings (id, kind, ref_id, content, vector, metadata, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, now()), now())
        ON CONFLICT (id) DO UPDATE
        SET kind = EXCLUDED.kind,
            ref_id = EXCLUDED.ref_id,
            content = EXCLUDED.content,
            vector = EXCLUDED.vector,
            metadata = EXCLUDED.metadata,
            updated_at = now()
    `, embedding.ID, embedding.Kind, embedding.RefID, embedding.Content, vector, embedding.Metadata, embedding.CreatedAt)
	return err
}

func (r *sqlRepository) DeleteByKind(ctx context.Context, kind string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM embeddings WHERE kind = $1`, kind)
	return err
}

func (r *sqlRepository) DeleteAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM embeddings WHERE vector IS NOT NULL`)
	return err
}

func (r *sqlRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM embeddings WHERE vector IS NOT NULL`); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *sqlRepository) Similar(ctx context.Context, vector []float32, limit int, minScore float64, kinds []string) ([]models.EmbeddingMatch, error) {
	if len(vector) != r.dimension {
		return nil, fmt.Errorf("embedding dimension mismatch: expected %d got %d", r.dimension, len(vector))
	}
	if limit <= 0 {
		limit = 5
	}
	args := []any{pgvector.NewVector(vector), limit}
	query := `
        SELECT id, kind, ref_id, content, metadata,
               1 - (vector <#> $1::vector) AS score
        FROM embeddings
        WHERE vector IS NOT NULL`
	if len(kinds) > 0 {
		query += " AND kind = ANY($3)"
		args = append(args, pq.Array(kinds))
	}
	query += " ORDER BY vector <#> $1::vector ASC LIMIT $2"

	rows := make([]models.EmbeddingMatch, 0, limit)
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	if minScore <= 0 {
		return rows, nil
	}
	filtered := rows[:0]
	for _, row := range rows {
		if row.Score >= minScore {
			filtered = append(filtered, row)
		}
	}
	return filtered, nil
}

func (r *sqlRepository) LoadConfig(ctx context.Context) (models.EmbeddingConfig, error) {
	var raw struct {
		Value []byte `db:"value"`
	}
	err := r.db.GetContext(ctx, &raw, `SELECT value FROM embedding_config WHERE key = $1`, personalizationConfigKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.EmbeddingConfig{Weight: 0.65}, nil
		}
		return models.EmbeddingConfig{}, err
	}
	cfg := models.EmbeddingConfig{Weight: 0.65}
	if len(raw.Value) > 0 {
		if err := json.Unmarshal(raw.Value, &cfg); err != nil {
			return models.EmbeddingConfig{}, fmt.Errorf("decode embedding config: %w", err)
		}
	}
	return cfg, nil
}

func (r *sqlRepository) SaveConfig(ctx context.Context, cfg models.EmbeddingConfig) error {
	payload, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode embedding config: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO embedding_config (key, value, updated_at)
        VALUES ($1, $2, now())
        ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()
    `, personalizationConfigKey, payload)
	return err
}
