package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	var items []domain.Category
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, name, created_at
		FROM resource_categories
		ORDER BY name ASC
	`)
	return items, err
}
