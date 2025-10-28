package repos

import (
    "context"
    "database/sql"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/tanydotai/tanyai/backend/internal/models"
)

// ExternalSourceRepository exposes DB operations for external_sources.
type ExternalSourceRepository interface {
    List(ctx context.Context, params ListParams) ([]models.ExternalSource, int64, error)
    Get(ctx context.Context, id uuid.UUID) (models.ExternalSource, error)
    FindByBaseURL(ctx context.Context, baseURL string) (models.ExternalSource, error)
    Create(ctx context.Context, source models.ExternalSource) (models.ExternalSource, error)
    Update(ctx context.Context, source models.ExternalSource) (models.ExternalSource, error)
    UpdateSyncState(ctx context.Context, id uuid.UUID, etag *string, lastModified *time.Time, syncedAt time.Time) error
    SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (models.ExternalSource, error)
    EnsureDefaults(ctx context.Context, defaults []models.ExternalSource) error
}

// NewExternalSourceRepository creates a SQL backed implementation.
func NewExternalSourceRepository(db *sqlx.DB) ExternalSourceRepository {
    return &externalSourceRepository{db: db}
}

type externalSourceRepository struct {
    db *sqlx.DB
}

func (r *externalSourceRepository) List(ctx context.Context, params ListParams) ([]models.ExternalSource, int64, error) {
    const baseQuery = `SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources`
    const countQuery = `SELECT COUNT(*) FROM external_sources`

    orderBy, err := params.ValidateSort(map[string]string{
        "name":       "LOWER(name)",
        "created_at": "created_at",
        "updated_at": "updated_at",
    }, "name")
    if err != nil {
        return nil, 0, err
    }

    query := baseQuery + " ORDER BY " + orderBy + " LIMIT $1 OFFSET $2"

    var sources []models.ExternalSource
    if err := r.db.SelectContext(ctx, &sources, query, params.Limit, params.Offset()); err != nil {
        return nil, 0, err
    }

    var total int64
    if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
        return nil, 0, err
    }

    return sources, total, nil
}

func (r *externalSourceRepository) Get(ctx context.Context, id uuid.UUID) (models.ExternalSource, error) {
    const query = `SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources WHERE id = $1`

    var source models.ExternalSource
    if err := r.db.GetContext(ctx, &source, query, id); err != nil {
        if err == sql.ErrNoRows {
            return models.ExternalSource{}, ErrNotFound
        }
        return models.ExternalSource{}, err
    }
    return source, nil
}

func (r *externalSourceRepository) FindByBaseURL(ctx context.Context, baseURL string) (models.ExternalSource, error) {
    const query = `SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources WHERE LOWER(base_url) = LOWER($1)`

    var source models.ExternalSource
    if err := r.db.GetContext(ctx, &source, query, strings.TrimSpace(baseURL)); err != nil {
        if err == sql.ErrNoRows {
            return models.ExternalSource{}, ErrNotFound
        }
        return models.ExternalSource{}, err
    }
    return source, nil
}

func (r *externalSourceRepository) Create(ctx context.Context, source models.ExternalSource) (models.ExternalSource, error) {
    const query = `INSERT INTO external_sources (name, base_url, source_type, enabled, etag, last_modified, last_synced_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at`

    var created models.ExternalSource
    if err := r.db.GetContext(ctx, &created, query,
        source.Name,
        source.BaseURL,
        source.SourceType,
        source.Enabled,
        source.ETag,
        source.LastModified,
        source.LastSyncedAt,
    ); err != nil {
        return models.ExternalSource{}, err
    }
    return created, nil
}

func (r *externalSourceRepository) Update(ctx context.Context, source models.ExternalSource) (models.ExternalSource, error) {
    const query = `UPDATE external_sources SET
    name = $2,
    base_url = $3,
    source_type = $4,
    enabled = $5,
    etag = $6,
    last_modified = $7,
    last_synced_at = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at`

    var updated models.ExternalSource
    if err := r.db.GetContext(ctx, &updated, query,
        source.ID,
        source.Name,
        source.BaseURL,
        source.SourceType,
        source.Enabled,
        source.ETag,
        source.LastModified,
        source.LastSyncedAt,
    ); err != nil {
        if err == sql.ErrNoRows {
            return models.ExternalSource{}, ErrNotFound
        }
        return models.ExternalSource{}, err
    }
    return updated, nil
}

func (r *externalSourceRepository) UpdateSyncState(ctx context.Context, id uuid.UUID, etag *string, lastModified *time.Time, syncedAt time.Time) error {
    const query = `UPDATE external_sources SET etag = $2, last_modified = $3, last_synced_at = $4, updated_at = NOW() WHERE id = $1`
    res, err := r.db.ExecContext(ctx, query, id, etag, lastModified, syncedAt)
    if err != nil {
        return err
    }
    affected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return ErrNotFound
    }
    return nil
}

func (r *externalSourceRepository) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (models.ExternalSource, error) {
    const query = `UPDATE external_sources SET enabled = $2, updated_at = NOW() WHERE id = $1 RETURNING id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at`
    var source models.ExternalSource
    if err := r.db.GetContext(ctx, &source, query, id, enabled); err != nil {
        if err == sql.ErrNoRows {
            return models.ExternalSource{}, ErrNotFound
        }
        return models.ExternalSource{}, err
    }
    return source, nil
}

func (r *externalSourceRepository) EnsureDefaults(ctx context.Context, defaults []models.ExternalSource) error {
    if len(defaults) == 0 {
        return nil
    }

    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        _ = tx.Rollback()
    }()

    for _, def := range defaults {
        if strings.TrimSpace(def.BaseURL) == "" || strings.TrimSpace(def.Name) == "" {
            continue
        }
        var existing models.ExternalSource
        err := tx.GetContext(ctx, &existing, `SELECT id, name, base_url, source_type, enabled, etag, last_modified, last_synced_at, created_at, updated_at FROM external_sources WHERE LOWER(base_url) = LOWER($1) LIMIT 1`, def.BaseURL)
        if err != nil {
            if err == sql.ErrNoRows {
                if _, err := tx.ExecContext(ctx, `INSERT INTO external_sources (name, base_url, source_type, enabled) VALUES ($1, $2, $3, $4)`, def.Name, def.BaseURL, def.SourceType, def.Enabled); err != nil {
                    return err
                }
                continue
            }
            return err
        }
        _, err = tx.ExecContext(ctx, `UPDATE external_sources SET name = $2, source_type = $3, enabled = $4, updated_at = NOW() WHERE id = $1`, existing.ID, def.Name, def.SourceType, def.Enabled)
        if err != nil {
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        return err
    }
    return nil
}
