package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
)

type ResourceRepo struct {
	db *sqlx.DB
}

func NewResourceRepo(db *sqlx.DB) *ResourceRepo {
	return &ResourceRepo{db: db}
}

func (r *ResourceRepo) List(ctx context.Context) ([]domain.Resource, error) {
	items := make([]domain.Resource, 0)
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, category_id, title, description, location, price_per_hour, is_active, created_at
		FROM resources
		ORDER BY id DESC
	`)
	return items, err
}

func (r *ResourceRepo) Create(
	ctx context.Context,
	categoryID uint64,
	title string,
	description, location *string,
	pricePerHour int,
) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO resources (category_id, title, description, location, price_per_hour)
		VALUES (?, ?, ?, ?, ?)
	`, categoryID, title, description, location, pricePerHour)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}
