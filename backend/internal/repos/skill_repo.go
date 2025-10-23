package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// SkillRepository exposes persistence operations for skills.
type SkillRepository interface {
	List(ctx context.Context, params ListParams) ([]models.Skill, int64, error)
	Create(ctx context.Context, skill models.Skill) (models.Skill, error)
	Update(ctx context.Context, skill models.Skill) (models.Skill, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, pairs []models.Skill) error
}

// NewSkillRepository constructs SkillRepository backed by SQL.
func NewSkillRepository(db *sqlx.DB) SkillRepository {
	return &skillRepository{db: db}
}

type skillRepository struct {
	db *sqlx.DB
}

func (r *skillRepository) List(ctx context.Context, params ListParams) ([]models.Skill, int64, error) {
	const baseQuery = `SELECT id, name, "order" FROM skills`
	const countQuery = `SELECT COUNT(*) FROM skills`

	orderBy, err := params.ValidateSort(map[string]string{
		"order": "\"order\"",
		"name":  "name",
	}, "order")
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery + " ORDER BY " + orderBy + " LIMIT $1 OFFSET $2"

	var skills []models.Skill
	if err := r.db.SelectContext(ctx, &skills, query, params.Limit, params.Offset()); err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	return skills, total, nil
}

func (r *skillRepository) Create(ctx context.Context, skill models.Skill) (models.Skill, error) {
	const query = `INSERT INTO skills (name, "order") VALUES ($1, $2) RETURNING id, name, "order"`

	var created models.Skill
	if err := r.db.GetContext(ctx, &created, query, skill.Name, skill.Order); err != nil {
		return models.Skill{}, err
	}
	return created, nil
}

func (r *skillRepository) Update(ctx context.Context, skill models.Skill) (models.Skill, error) {
	const query = `UPDATE skills SET name = $2, "order" = $3 WHERE id = $1 RETURNING id, name, "order"`

	var updated models.Skill
	if err := r.db.GetContext(ctx, &updated, query, skill.ID, skill.Name, skill.Order); err != nil {
		if err == sql.ErrNoRows {
			return models.Skill{}, ErrNotFound
		}
		return models.Skill{}, err
	}
	return updated, nil
}

func (r *skillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM skills WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *skillRepository) Reorder(ctx context.Context, pairs []models.Skill) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, pair := range pairs {
		res, err := tx.ExecContext(ctx, `UPDATE skills SET "order" = $1 WHERE id = $2`, pair.Order, pair.ID)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrNotFound
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
