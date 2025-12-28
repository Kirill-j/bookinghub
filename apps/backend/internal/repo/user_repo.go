package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"

	"database/sql"
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

func (r *UserRepo) UpdateProfile(ctx context.Context, id uint64, email, name string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET email = ?, name = ?
		WHERE id = ?
	`, email, name, id)
	return err
}

func (r *UserRepo) UpdatePasswordHashByID(ctx context.Context, id uint64, hash string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = ?
		WHERE id = ?
	`, hash, id)
	return err
}

func (r *UserRepo) DeleteAccount(ctx context.Context, userID uint64) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Удаляем брони, которые сделал этот пользователь
	if _, err := tx.ExecContext(ctx, `DELETE FROM bookings WHERE user_id = ?`, userID); err != nil {
		return err
	}

	// 2) Удаляем брони по ресурсам, которыми владеет пользователь
	if _, err := tx.ExecContext(ctx, `
		DELETE b FROM bookings b
		JOIN resources r ON r.id = b.resource_id
		WHERE r.owner_user_id = ?
	`, userID); err != nil {
		return err
	}

	// 3) Удаляем ресурсы пользователя (объявления)
	if _, err := tx.ExecContext(ctx, `DELETE FROM resources WHERE owner_user_id = ?`, userID); err != nil {
		return err
	}

	// 4) Удаляем самого пользователя
	if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepo) GetRoleByID(ctx context.Context, id uint64) (domain.UserRole, error) {
	var role domain.UserRole
	err := r.db.GetContext(ctx, &role, `
		SELECT role
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)
	return role, err
}
