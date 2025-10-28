package repos

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/tanydotai/tanyai/backend/internal/models"
)

// ExternalItemWithSource couples an external item with its source metadata.
type ExternalItemWithSource struct {
    models.ExternalItem
    SourceName          string `db:"source_name"`
    SourceBaseURL       string `db:"source_base_url"`
}

// ExternalItemListParams controls listing behaviour.
type ExternalItemListParams struct {
    ListParams
    SourceID *uuid.UUID
    Kind     string
    Visible  *bool
    Search   string
}

// ExternalItemRepository defines persistence operations for external_items.
type ExternalItemRepository interface {
    Upsert(ctx context.Context, items []models.ExternalItem) error
    List(ctx context.Context, params ExternalItemListParams) ([]ExternalItemWithSource, int64, error)
    SetVisibility(ctx context.Context, id uuid.UUID, visible bool) (ExternalItemWithSource, error)
}

// NewExternalItemRepository constructs a repository instance.
func NewExternalItemRepository(db *sqlx.DB) ExternalItemRepository {
    return &externalItemRepository{db: db}
}

type externalItemRepository struct {
    db *sqlx.DB
}

func (r *externalItemRepository) Upsert(ctx context.Context, items []models.ExternalItem) error {
    if len(items) == 0 {
        return nil
    }

    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        _ = tx.Rollback()
    }()

    const query = `INSERT INTO external_items (source_id, kind, title, url, summary, content, metadata, published_at, hash, visible)
VALUES (:source_id, :kind, :title, :url, :summary, :content, :metadata, :published_at, :hash, :visible)
ON CONFLICT (source_id, hash) DO UPDATE SET
    title = EXCLUDED.title,
    url = EXCLUDED.url,
    summary = EXCLUDED.summary,
    content = EXCLUDED.content,
    metadata = EXCLUDED.metadata,
    published_at = EXCLUDED.published_at,
    visible = EXCLUDED.visible,
    updated_at = NOW()`

    for _, item := range items {
        if item.Metadata == nil {
            item.Metadata = models.JSONB{}
        }
        if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        return err
    }
    return nil
}

func (r *externalItemRepository) List(ctx context.Context, params ExternalItemListParams) ([]ExternalItemWithSource, int64, error) {
    builder := strings.Builder{}
    builder.WriteString(`SELECT i.id, i.source_id, i.kind, i.title, i.url, i.summary, i.content, i.metadata, i.published_at, i.hash, i.visible, i.created_at, i.updated_at, s.name AS source_name, s.base_url AS source_base_url
FROM external_items i
JOIN external_sources s ON s.id = i.source_id`)

    var where []string
    args := make([]any, 0, 5)

    if params.SourceID != nil {
        where = append(where, fmt.Sprintf("i.source_id = $%d", len(args)+1))
        args = append(args, *params.SourceID)
    }
    if params.Kind != "" {
        where = append(where, fmt.Sprintf("LOWER(i.kind) = LOWER($%d)", len(args)+1))
        args = append(args, params.Kind)
    }
    if params.Visible != nil {
        where = append(where, fmt.Sprintf("i.visible = $%d", len(args)+1))
        args = append(args, *params.Visible)
    }
    if strings.TrimSpace(params.Search) != "" {
        search := "%%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%%"
        where = append(where, fmt.Sprintf("(LOWER(i.title) LIKE $%d OR LOWER(i.summary) LIKE $%d)", len(args)+1, len(args)+2))
        args = append(args, search, search)
    }

    if len(where) > 0 {
        builder.WriteString(" WHERE ")
        builder.WriteString(strings.Join(where, " AND "))
    }

    sortParams := params
    if sortParams.SortField == "" && sortParams.SortDir == "" {
        sortParams.SortDir = "desc"
    }

    orderBy, err := sortParams.ValidateSort(map[string]string{
        "published_at": "i.published_at",
        "updated_at":   "i.updated_at",
        "title":        "LOWER(i.title)",
    }, "published_at")
    if err != nil {
        return nil, 0, err
    }
    filterArgs := append([]any(nil), args...)

    builder.WriteString(" ORDER BY ")
    builder.WriteString(orderBy)
    builder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2))
    args = append(args, params.Limit, params.Offset())

    query := builder.String()
    rows := make([]ExternalItemWithSource, 0, params.Limit)
    if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
        return nil, 0, err
    }

    countBuilder := strings.Builder{}
    countBuilder.WriteString("SELECT COUNT(*) FROM external_items i")
    if len(where) > 0 {
        countBuilder.WriteString(" WHERE ")
        countBuilder.WriteString(strings.Join(where, " AND "))
    }
    countQuery := countBuilder.String()
    var total int64
    if err := r.db.GetContext(ctx, &total, countQuery, filterArgs...); err != nil {
        return nil, 0, err
    }

    return rows, total, nil
}

func (r *externalItemRepository) SetVisibility(ctx context.Context, id uuid.UUID, visible bool) (ExternalItemWithSource, error) {
    const query = `UPDATE external_items SET visible = $2, updated_at = NOW() WHERE id = $1 RETURNING id, source_id, kind, title, url, summary, content, metadata, published_at, hash, visible, created_at, updated_at, (SELECT name FROM external_sources WHERE id = external_items.source_id) AS source_name, (SELECT base_url FROM external_sources WHERE id = external_items.source_id) AS source_base_url`
    var item ExternalItemWithSource
    if err := r.db.GetContext(ctx, &item, query, id, visible); err != nil {
        if err == sql.ErrNoRows {
            return ExternalItemWithSource{}, ErrNotFound
        }
        return ExternalItemWithSource{}, err
    }
    return item, nil
}
