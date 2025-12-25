package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
)

type BookingRepo struct {
	db *sqlx.DB
}

func NewBookingRepo(db *sqlx.DB) *BookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) ListByUser(ctx context.Context, userID uint64) ([]domain.Booking, error) {
	var items []domain.Booking
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE user_id = ?
		ORDER BY start_at DESC
	`, userID)
	return items, err
}

func (r *BookingRepo) ListPending(ctx context.Context) ([]domain.Booking, error) {
	var items []domain.Booking
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE status = 'PENDING'
		ORDER BY start_at ASC
	`)
	return items, err
}

func (r *BookingRepo) Create(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO bookings (resource_id, user_id, start_at, end_at, status)
		VALUES (?, ?, ?, ?, 'PENDING')
	`, resourceID, userID, startAt, endAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *BookingRepo) HasConflict(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
	// Пересечение: start < existing_end AND end > existing_start
	// Считаем конфликтами PENDING и APPROVED
	var cnt int
	err := r.db.GetContext(ctx, &cnt, `
		SELECT COUNT(*)
		FROM bookings
		WHERE resource_id = ?
		  AND status IN ('PENDING','APPROVED')
		  AND (? < end_at) AND (? > start_at)
	`, resourceID, startAt, endAt)

	return cnt > 0, err
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id uint64, status domain.BookingStatus, managerComment *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE bookings
		SET status = ?, manager_comment = ?
		WHERE id = ?
	`, status, managerComment, id)
	return err
}

func (r *BookingRepo) GetByID(ctx context.Context, id uint64) (*domain.Booking, error) {
	var b domain.Booking
	err := r.db.GetContext(ctx, &b, `
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE id = ?
		LIMIT 1
	`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepo) Cancel(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE bookings
		SET status = 'CANCELED'
		WHERE id = ?
	`, id)
	return err
}

func (r *BookingRepo) ListByResourceBetween(ctx context.Context, resourceID uint64, from, to time.Time) ([]domain.Booking, error) {
	var items []domain.Booking
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE resource_id = ?
		  AND status IN ('PENDING','APPROVED')
		  AND start_at >= ?
		  AND start_at < ?
		ORDER BY start_at ASC
	`, resourceID, from, to)
	return items, err
}
