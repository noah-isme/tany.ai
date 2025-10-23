package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ServiceRepository defines CRUD behaviour for services.
type ServiceRepository interface {
	List(ctx context.Context, params ListParams) ([]models.Service, int64, error)
	Create(ctx context.Context, service models.Service) (models.Service, error)
	Update(ctx context.Context, service models.Service) (models.Service, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, pairs []models.Service) error
	Toggle(ctx context.Context, id uuid.UUID, desired *bool) (models.Service, error)
}

// NewServiceRepository constructs SQL-backed service repository.
func NewServiceRepository(db *sqlx.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

type serviceRepository struct {
	db *sqlx.DB
}

func (r *serviceRepository) List(ctx context.Context, params ListParams) ([]models.Service, int64, error) {
	const baseQuery = `SELECT id, name, description, price_min, price_max, currency, duration_label, is_active, "order" FROM services`
	const countQuery = `SELECT COUNT(*) FROM services`

	orderBy, err := params.ValidateSort(map[string]string{
		"order": "\"order\"",
		"name":  "name",
	}, "order")
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery + " ORDER BY " + orderBy + " LIMIT $1 OFFSET $2"

	var services []models.Service
	if err := r.db.SelectContext(ctx, &services, query, params.Limit, params.Offset()); err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepository) Create(ctx context.Context, service models.Service) (models.Service, error) {
	const query = `INSERT INTO services (name, description, price_min, price_max, currency, duration_label, is_active, "order")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, name, description, price_min, price_max, currency, duration_label, is_active, "order"`

	var created models.Service
	if err := r.db.GetContext(ctx, &created, query,
		service.Name,
		service.Description,
		service.PriceMin,
		service.PriceMax,
		service.Currency,
		service.DurationLabel,
		service.IsActive,
		service.Order,
	); err != nil {
		return models.Service{}, err
	}
	return created, nil
}

func (r *serviceRepository) Update(ctx context.Context, service models.Service) (models.Service, error) {
	const query = `UPDATE services SET
    name = $2,
    description = $3,
    price_min = $4,
    price_max = $5,
    currency = $6,
    duration_label = $7,
    is_active = $8,
    "order" = $9
WHERE id = $1
RETURNING id, name, description, price_min, price_max, currency, duration_label, is_active, "order"`

	var updated models.Service
	if err := r.db.GetContext(ctx, &updated, query,
		service.ID,
		service.Name,
		service.Description,
		service.PriceMin,
		service.PriceMax,
		service.Currency,
		service.DurationLabel,
		service.IsActive,
		service.Order,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Service{}, ErrNotFound
		}
		return models.Service{}, err
	}
	return updated, nil
}

func (r *serviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM services WHERE id = $1`
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

func (r *serviceRepository) Reorder(ctx context.Context, pairs []models.Service) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, pair := range pairs {
		res, err := tx.ExecContext(ctx, `UPDATE services SET "order" = $1 WHERE id = $2`, pair.Order, pair.ID)
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

func (r *serviceRepository) Toggle(ctx context.Context, id uuid.UUID, desired *bool) (models.Service, error) {
	query := `UPDATE services SET is_active = NOT is_active WHERE id = $1 RETURNING id, name, description, price_min, price_max, currency, duration_label, is_active, "order"`
	args := []any{id}
	if desired != nil {
		query = `UPDATE services SET is_active = $2 WHERE id = $1 RETURNING id, name, description, price_min, price_max, currency, duration_label, is_active, "order"`
		args = append(args, *desired)
	}

	var service models.Service
	if err := r.db.GetContext(ctx, &service, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return models.Service{}, ErrNotFound
		}
		return models.Service{}, err
	}
	return service, nil
}
