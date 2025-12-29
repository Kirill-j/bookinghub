package repo

import (
	"context"
	"testing"
	"time"

	"bookinghub-backend/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newSQLXMockRepo3(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func TestBookingRepo_UpdateStatus_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo3(t)
	defer closeFn()

	r := NewBookingRepo(db)

	comment := "ok"
	mock.ExpectExec("UPDATE bookings\\s+SET status = \\?, manager_comment = \\?\\s+WHERE id = \\?").
		WithArgs("APPROVED", &comment, uint64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.UpdateStatus(context.Background(), 10, domain.BookingApproved, &comment); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookingRepo_ListByResourceBetween_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo3(t)
	defer closeFn()

	r := NewBookingRepo(db)

	from := time.Now().Add(24 * time.Hour)
	to := from.Add(48 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "resource_id", "user_id", "start_at", "end_at", "status", "manager_comment", "created_at", "updated_at",
	}).AddRow(
		uint64(1), uint64(7), uint64(99),
		from.Add(2*time.Hour), from.Add(3*time.Hour),
		"PENDING", nil, time.Now(), nil,
	)

	mock.ExpectQuery("FROM bookings\\s+WHERE resource_id = \\?").
		WithArgs(uint64(7), from, to).
		WillReturnRows(rows)

	items, err := r.ListByResourceBetween(context.Background(), 7, from, to)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
