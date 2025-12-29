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

func (r *CategoryRepo) Create(ctx context.Context, name string) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO resource_categories (name) VALUES (?)
	`, name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *CategoryRepo) Update(ctx context.Context, id uint64, name string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE resource_categories
		SET name = ?
		WHERE id = ?
	`, name, id)
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM resource_categories WHERE id = ?`, id)
	return err
}
