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

func TestExternalItemRepositoryUpsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalItemRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO external_items (source_id, kind, title, url, summary, content, metadata, published_at, hash, visible)`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id := uuid.New()
	items := []models.ExternalItem{{
		SourceID: id,
		Kind:     "post",
		Title:    "Hello",
		URL:      "https://noahis.me/post",
		Metadata: models.JSONB{"sourceName": "noahis.me"},
		Visible:  true,
		Hash:     "hash",
	}}

	err = repo.Upsert(context.Background(), items)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExternalItemRepositoryListWithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalItemRepository(sqlx.NewDb(db, "sqlmock"))

	sourceID := uuid.New()
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "source_id", "kind", "title", "url", "summary", "content", "metadata", "published_at", "hash", "visible", "created_at", "updated_at", "source_name", "source_base_url"}).
		AddRow(uuid.New(), sourceID, "post", "Hello", "https://noahis.me/post", "Summary", nil, []byte(`{"sourceName":"noahis.me"}`), now, "hash", true, now, now, "noahis.me", "https://noahis.me")

	search := "%%golang%%"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT i.id, i.source_id, i.kind, i.title, i.url, i.summary, i.content, i.metadata, i.published_at, i.hash, i.visible, i.created_at, i.updated_at, s.name AS source_name, s.base_url AS source_base_url
FROM external_items i
JOIN external_sources s ON s.id = i.source_id WHERE i.source_id = $1 AND LOWER(i.kind) = LOWER($2) AND i.visible = $3 AND (LOWER(i.title) LIKE $4 OR LOWER(i.summary) LIKE $5) ORDER BY i.published_at DESC LIMIT $6 OFFSET $7`)).
		WithArgs(sourceID, "post", true, search, search, 20, 0).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM external_items i WHERE i.source_id = $1 AND LOWER(i.kind) = LOWER($2) AND i.visible = $3 AND (LOWER(i.title) LIKE $4 OR LOWER(i.summary) LIKE $5)`)).
		WithArgs(sourceID, "post", true, search, search).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	visible := true
	params := ExternalItemListParams{
		ListParams: ListParams{Page: 1, Limit: 20, SortField: "published_at", SortDir: "desc"},
		SourceID:   &sourceID,
		Kind:       "post",
		Visible:    &visible,
		Search:     "golang",
	}

	items, total, err := repo.List(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, items, 1)
	require.Equal(t, "Hello", items[0].Title)
	require.Equal(t, "noahis.me", items[0].SourceName)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExternalItemRepositorySetVisibilityNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewExternalItemRepository(sqlx.NewDb(db, "sqlmock"))

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE external_items SET visible = $2, updated_at = NOW() WHERE id = $1 RETURNING id, source_id, kind, title, url, summary, content, metadata, published_at, hash, visible, created_at, updated_at, (SELECT name FROM external_sources WHERE id = external_items.source_id) AS source_name, (SELECT base_url FROM external_sources WHERE id = external_items.source_id) AS source_base_url`)).
		WithArgs(id, false).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.SetVisibility(context.Background(), id, false)
	require.ErrorIs(t, err, ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}
