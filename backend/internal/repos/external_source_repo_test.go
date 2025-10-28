package repos

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

func TestExternalSourceRepositoryList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalSourceRepository(sqlx.NewDb(db, "sqlmock"))

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "base_url", "source_type", "enabled", "etag", "last_modified", "last_synced_at", "created_at", "updated_at"}).
		AddRow(uuid.New(), "noahis.me", "https://noahis.me", "auto", true, nil, now, now, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources ORDER BY LOWER(name) ASC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM external_sources`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	sources, total, err := repo.List(context.Background(), ListParams{Page: 1, Limit: 10, SortField: "name", SortDir: "asc"})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, sources, 1)
	require.Equal(t, "noahis.me", sources[0].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExternalSourceRepositoryUpdateSyncStateNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalSourceRepository(sqlx.NewDb(db, "sqlmock"))

	id := uuid.New()
	now := time.Now()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE external_sources SET etag = $2, last_modified = $3, last_synced_at = $4, updated_at = NOW() WHERE id = $1`)).
		WithArgs(id, nil, nil, now).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateSyncState(context.Background(), id, nil, nil, now)
	require.ErrorIs(t, err, ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExternalSourceRepositoryEnsureDefaultsInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalSourceRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources WHERE LOWER(base_url) = LOWER($1) LIMIT 1`)).
		WithArgs("https://noahis.me").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO external_sources (name, base_url, source_type, enabled) VALUES ($1, $2, $3, $4)`)).
		WithArgs("noahis.me", "https://noahis.me", "auto", true).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	defaults := []models.ExternalSource{{
		Name:       "noahis.me",
		BaseURL:    "https://noahis.me",
		SourceType: "auto",
		Enabled:    true,
	}}

	err = repo.EnsureDefaults(context.Background(), defaults)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
