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

// UserRepository exposes persistence operations for users and refresh tokens.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (models.User, []string, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.User, []string, error)
	CreateRefreshToken(ctx context.Context, token models.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, hash string) (models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) (bool, error)
	RevokeRefreshTokenByHash(ctx context.Context, hash string) error
	DeleteExpiredRefreshTokens(ctx context.Context, before time.Time) (int64, error)
}

// NewUserRepository constructs a SQL-backed user repository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (models.User, []string, error) {
	const query = `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE LOWER(email) = LOWER($1)`

	var user models.User
	if err := r.db.GetContext(ctx, &user, query, strings.TrimSpace(email)); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil, ErrNotFound
		}
		return models.User{}, nil, err
	}

	roles, err := r.fetchRoles(ctx, user.ID)
	if err != nil {
		return models.User{}, nil, err
	}
	return user, roles, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (models.User, []string, error) {
	const query = `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`
	var user models.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil, ErrNotFound
		}
		return models.User{}, nil, err
	}

	roles, err := r.fetchRoles(ctx, user.ID)
	if err != nil {
		return models.User{}, nil, err
	}
	return user, roles, nil
}

func (r *userRepository) fetchRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	const query = `SELECT role FROM user_roles WHERE user_id = $1 ORDER BY role`
	var roles []string
	if err := r.db.SelectContext(ctx, &roles, query, userID); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, token models.RefreshToken) error {
	const query = `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.Revoked,
	)
	return err
}

func (r *userRepository) FindRefreshTokenByHash(ctx context.Context, hash string) (models.RefreshToken, error) {
	const query = `SELECT id, user_id, token_hash, expires_at, revoked, created_at FROM refresh_tokens WHERE token_hash = $1`
	var token models.RefreshToken
	if err := r.db.GetContext(ctx, &token, query, hash); err != nil {
		if err == sql.ErrNoRows {
			return models.RefreshToken{}, ErrNotFound
		}
		return models.RefreshToken{}, err
	}
	return token, nil
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1 AND revoked = FALSE`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (r *userRepository) RevokeRefreshTokenByHash(ctx context.Context, hash string) error {
	const query = `UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1`
	_, err := r.db.ExecContext(ctx, query, hash)
	return err
}

func (r *userRepository) DeleteExpiredRefreshTokens(ctx context.Context, before time.Time) (int64, error) {
	const query = `DELETE FROM refresh_tokens WHERE expires_at <= $1`
	res, err := r.db.ExecContext(ctx, query, before)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
