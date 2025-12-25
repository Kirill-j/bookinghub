package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.GetContext(ctx, &u, `
		SELECT id, email, name, role, password_hash, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, email, name string, role domain.UserRole, passwordHash string) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO users (email, name, role, password_hash)
		VALUES (?, ?, ?, ?)
	`, email, name, role, passwordHash)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *UserRepo) UpdatePasswordHash(ctx context.Context, email, hash string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET password_hash = ?
		WHERE email = ?
	`, hash, email)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	var u domain.User
	err := r.db.GetContext(ctx, &u, `
		SELECT id, email, name, role, password_hash, created_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
