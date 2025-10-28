package embedding

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

func TestRepositoryUpsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(sqlx.NewDb(db, "sqlmock"), 1536)
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO embeddings (id, kind, ref_id, content, vector, metadata, created_at, updated_at)`)).
		WithArgs(id, "profile", nil, "Persona profesional.", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Upsert(context.Background(), models.Embedding{
		ID:      id,
		Kind:    "profile",
		Content: "Persona profesional.",
		Vector:  make([]float32, 1536),
		Metadata: models.JSONB{
			"name": "Noah",
		},
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositorySimilarFiltersByScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxdb, 3)

	rows := sqlmock.NewRows([]string{"id", "kind", "ref_id", "content", "metadata", "score"}).
		AddRow(uuid.New(), "profile", nil, "Persona profesional.", `{"name":"Noah"}`, 0.9).
		AddRow(uuid.New(), "service", nil, "Layanan A", `{"name":"A"}`, 0.3)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, kind, ref_id, content, metadata,
               1 - (vector <#> $1::vector) AS score
        FROM embeddings
        WHERE vector IS NOT NULL ORDER BY vector <#> $1::vector ASC LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), 3).
		WillReturnRows(rows)

	matches, err := repo.Similar(context.Background(), []float32{1, 0, 1}, 3, 0.5, nil)
	require.NoError(t, err)
	require.Len(t, matches, 1)
	require.Equal(t, "profile", matches[0].Kind)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryLoadConfigDefault(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxdb, 1536)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM embedding_config WHERE key = $1`)).
		WithArgs(personalizationConfigKey).
		WillReturnError(sql.ErrNoRows)

	cfg, err := repo.LoadConfig(context.Background())
	require.NoError(t, err)
	require.InDelta(t, 0.65, cfg.Weight, 1e-9)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositorySaveConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxdb, 1536)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO embedding_config (key, value, updated_at)
        VALUES ($1, $2, now())
        ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`)).
		WithArgs(personalizationConfigKey, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveConfig(context.Background(), models.EmbeddingConfig{Weight: 0.7})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
