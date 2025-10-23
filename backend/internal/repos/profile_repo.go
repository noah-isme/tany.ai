package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ProfileRepository defines data access behaviour for profile entity.
type ProfileRepository interface {
	Get(ctx context.Context) (models.Profile, error)
	Upsert(ctx context.Context, profile models.Profile) (models.Profile, error)
}

// NewProfileRepository constructs a ProfileRepository backed by SQL.
func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &profileRepository{db: db}
}

type profileRepository struct {
	db *sqlx.DB
}

func (r *profileRepository) Get(ctx context.Context) (models.Profile, error) {
	const query = `SELECT id, name, title, bio, email, phone, location, avatar_url, updated_at FROM profile LIMIT 1`

	var profile models.Profile
	if err := r.db.GetContext(ctx, &profile, query); err != nil {
		if err == sql.ErrNoRows {
			return models.Profile{}, ErrNotFound
		}
		return models.Profile{}, err
	}
	return profile, nil
}

func (r *profileRepository) Upsert(ctx context.Context, profile models.Profile) (models.Profile, error) {
	if profile.ID == uuid.Nil {
		existing, err := r.Get(ctx)
		if err == nil {
			profile.ID = existing.ID
		} else if err != ErrNotFound {
			return models.Profile{}, err
		}
	}

	if profile.ID == uuid.Nil {
		const insertQuery = `INSERT INTO profile (name, title, bio, email, phone, location, avatar_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, title, bio, email, phone, location, avatar_url, updated_at`
		var created models.Profile
		if err := r.db.GetContext(ctx, &created, insertQuery,
			profile.Name,
			profile.Title,
			profile.Bio,
			profile.Email,
			profile.Phone,
			profile.Location,
			profile.AvatarURL,
		); err != nil {
			return models.Profile{}, err
		}
		return created, nil
	}

	const upsertQuery = `INSERT INTO profile (id, name, title, bio, email, phone, location, avatar_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    title = EXCLUDED.title,
    bio = EXCLUDED.bio,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    location = EXCLUDED.location,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = NOW()
RETURNING id, name, title, bio, email, phone, location, avatar_url, updated_at`

	var updated models.Profile
	if err := r.db.GetContext(ctx, &updated, upsertQuery,
		profile.ID,
		profile.Name,
		profile.Title,
		profile.Bio,
		profile.Email,
		profile.Phone,
		profile.Location,
		profile.AvatarURL,
	); err != nil {
		return models.Profile{}, err
	}
	return updated, nil
}
