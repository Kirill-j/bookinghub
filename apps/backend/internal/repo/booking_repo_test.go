package repo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBookingRepo_ListByUser(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)
	now := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id", "resource_id", "user_id", "start_at", "end_at", "status",
		"manager_comment", "created_at", "updated_at",
	}).AddRow(
		uint64(1), uint64(10), uint64(5),
		now, now.Add(time.Hour), "PENDING",
		nil, now, nil,
	)

	q := regexp.QuoteMeta(`
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE user_id = ?
		ORDER BY start_at DESC
	`)
	mock.ExpectQuery(q).WithArgs(uint64(5)).WillReturnRows(rows)

	items, err := r.ListByUser(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListByUser err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != 1 || items[0].ResourceID != 10 || items[0].UserID != 5 {
		t.Fatalf("unexpected item: %+v", items[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_ListPending(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)
	now := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id", "resource_id", "user_id", "start_at", "end_at", "status",
		"manager_comment", "created_at", "updated_at",
	}).AddRow(uint64(2), uint64(11), uint64(6), now, now.Add(time.Hour), "PENDING", nil, now, nil)

	q := regexp.QuoteMeta(`
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE status = 'PENDING'
		ORDER BY start_at ASC
	`)
	mock.ExpectQuery(q).WillReturnRows(rows)

	items, err := r.ListPending(context.Background())
	if err != nil {
		t.Fatalf("ListPending err: %v", err)
	}
	if len(items) != 1 || items[0].ID != 2 {
		t.Fatalf("unexpected items: %+v", items)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_Create(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)
	start := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	q := regexp.QuoteMeta(`
		INSERT INTO bookings (resource_id, user_id, start_at, end_at, status)
		VALUES (?, ?, ?, ?, 'PENDING')
	`)

	mock.ExpectExec(q).
		WithArgs(uint64(7), uint64(9), start, end).
		WillReturnResult(sqlmock.NewResult(123, 1))

	id, err := r.Create(context.Background(), 7, 9, start, end)
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}
	if id != 123 {
		t.Fatalf("expected id=123 got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_HasConflict(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)
	start := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	q := regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM bookings
		WHERE resource_id = ?
		  AND status IN ('PENDING','APPROVED')
		  AND (? < end_at) AND (? > start_at)
	`)

	// конфликт есть
	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1)
	mock.ExpectQuery(q).WithArgs(uint64(99), start, end).WillReturnRows(rows)

	ok, err := r.HasConflict(context.Background(), 99, start, end)
	if err != nil {
		t.Fatalf("HasConflict err: %v", err)
	}
	if !ok {
		t.Fatalf("expected conflict=true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_GetByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)

	q := regexp.QuoteMeta(`
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE id = ?
		LIMIT 1
	`)
	mock.ExpectQuery(q).WithArgs(uint64(777)).WillReturnError(sql.ErrNoRows)

	b, err := r.GetByID(context.Background(), 777)
	if err != nil {
		t.Fatalf("GetByID err: %v", err)
	}
	if b != nil {
		t.Fatalf("expected nil, got %+v", b)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_Cancel(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)

	q := regexp.QuoteMeta(`
		UPDATE bookings
		SET status = 'CANCELED'
		WHERE id = ?
	`)
	mock.ExpectExec(q).WithArgs(uint64(55)).WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.Cancel(context.Background(), 55)
	if err != nil {
		t.Fatalf("Cancel err: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_GetOwnerUserIDByBookingID(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)

	q := regexp.QuoteMeta(`
		SELECT r.owner_user_id
		FROM bookings b
		JOIN resources r ON r.id = b.resource_id
		WHERE b.id = ?
	`)
	rows := sqlmock.NewRows([]string{"owner_user_id"}).AddRow(uint64(99))
	mock.ExpectQuery(q).WithArgs(uint64(1)).WillReturnRows(rows)

	owner, err := r.GetOwnerUserIDByBookingID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetOwnerUserIDByBookingID err: %v", err)
	}
	if owner != 99 {
		t.Fatalf("expected 99, got %d", owner)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingRepo_ListPendingForOwner(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewBookingRepo(db)
	now := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	q := regexp.QuoteMeta(`
		SELECT b.id, b.resource_id, b.user_id, b.start_at, b.end_at, b.status, b.manager_comment, b.created_at, b.updated_at
		FROM bookings b
		JOIN resources r ON r.id = b.resource_id
		WHERE b.status = 'PENDING'
		  AND r.owner_user_id = ?
		ORDER BY b.start_at ASC
	`)

	rows := sqlmock.NewRows([]string{
		"id", "resource_id", "user_id", "start_at", "end_at", "status",
		"manager_comment", "created_at", "updated_at",
	}).AddRow(uint64(1), uint64(10), uint64(3), now, now.Add(time.Hour), "PENDING", nil, now, nil)

	mock.ExpectQuery(q).WithArgs(uint64(42)).WillReturnRows(rows)

	items, err := r.ListPendingForOwner(context.Background(), 42)
	if err != nil {
		t.Fatalf("ListPendingForOwner err: %v", err)
	}
	if len(items) != 1 || items[0].UserID != 3 {
		t.Fatalf("unexpected items: %+v", items)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
