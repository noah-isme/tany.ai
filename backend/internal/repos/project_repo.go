package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ProjectRepository defines DB operations for projects.
type ProjectRepository interface {
	List(ctx context.Context, params ListParams) ([]models.Project, int64, error)
	Create(ctx context.Context, project models.Project) (models.Project, error)
	Update(ctx context.Context, project models.Project) (models.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, pairs []models.Project) error
	SetFeatured(ctx context.Context, id uuid.UUID, featured bool) (models.Project, error)
}

// NewProjectRepository constructs SQL-backed project repository.
func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepository{db: db}
}

type projectRepository struct {
	db *sqlx.DB
}

func (r *projectRepository) List(ctx context.Context, params ListParams) ([]models.Project, int64, error) {
	const baseQuery = `SELECT id, title, description, tech_stack, image_url, project_url, category, "order", is_featured FROM projects`
	const countQuery = `SELECT COUNT(*) FROM projects`

	orderBy, err := params.ValidateSort(map[string]string{
		"order":       "\"order\"",
		"title":       "title",
		"is_featured": "is_featured",
	}, "order")
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery + " ORDER BY " + orderBy + " LIMIT $1 OFFSET $2"

	var projects []models.Project
	if err := r.db.SelectContext(ctx, &projects, query, params.Limit, params.Offset()); err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) Create(ctx context.Context, project models.Project) (models.Project, error) {
	const query = `INSERT INTO projects (title, description, tech_stack, image_url, project_url, category, "order", is_featured)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, description, tech_stack, image_url, project_url, category, "order", is_featured`

	var created models.Project
	if err := r.db.GetContext(ctx, &created, query,
		project.Title,
		project.Description,
		project.TechStack,
		project.ImageURL,
		project.ProjectURL,
		project.Category,
		project.Order,
		project.IsFeatured,
	); err != nil {
		return models.Project{}, err
	}
	return created, nil
}

func (r *projectRepository) Update(ctx context.Context, project models.Project) (models.Project, error) {
	const query = `UPDATE projects SET
    title = $2,
    description = $3,
    tech_stack = $4,
    image_url = $5,
    project_url = $6,
    category = $7,
    "order" = $8,
    is_featured = $9
WHERE id = $1
RETURNING id, title, description, tech_stack, image_url, project_url, category, "order", is_featured`

	var updated models.Project
	if err := r.db.GetContext(ctx, &updated, query,
		project.ID,
		project.Title,
		project.Description,
		project.TechStack,
		project.ImageURL,
		project.ProjectURL,
		project.Category,
		project.Order,
		project.IsFeatured,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, err
	}
	return updated, nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM projects WHERE id = $1`
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

func (r *projectRepository) Reorder(ctx context.Context, pairs []models.Project) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, pair := range pairs {
		res, err := tx.ExecContext(ctx, `UPDATE projects SET "order" = $1 WHERE id = $2`, pair.Order, pair.ID)
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

func (r *projectRepository) SetFeatured(ctx context.Context, id uuid.UUID, featured bool) (models.Project, error) {
	const query = `UPDATE projects SET is_featured = $2 WHERE id = $1 RETURNING id, title, description, tech_stack, image_url, project_url, category, "order", is_featured`

	var project models.Project
	if err := r.db.GetContext(ctx, &project, query, id, featured); err != nil {
		if err == sql.ErrNoRows {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, err
	}
	return project, nil
}
