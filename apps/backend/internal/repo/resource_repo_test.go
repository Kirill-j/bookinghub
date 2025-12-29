package repo

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newRepoMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock, func() { _ = db.Close() }
}

func TestResourceRepo_List_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock(t)
	defer cleanup()

	r := NewResourceRepo(dbx)
	now := time.Now()

	mock.ExpectQuery("SELECT id, owner_user_id, category_id, title, description, location, price_per_hour, is_active, created_at FROM resources ORDER BY id DESC").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "owner_user_id", "category_id", "title", "description", "location", "price_per_hour", "is_active", "created_at",
		}).AddRow(uint64(1), uint64(2), uint64(3), "T", nil, nil, 10, true, now))

	_, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestResourceRepo_Create_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock(t)
	defer cleanup()

	r := NewResourceRepo(dbx)

	mock.ExpectExec("INSERT INTO resources \\(owner_user_id, category_id, title, description, location, price_per_hour\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(uint64(2), uint64(3), "T", nil, nil, 10).
		WillReturnResult(sqlmock.NewResult(5, 1))

	id, err := r.Create(context.Background(), 2, 3, "T", nil, nil, 10)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if id != 5 {
		t.Fatalf("expected id=5 got %d", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestResourceRepo_ListByOwner_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock(t)
	defer cleanup()

	r := NewResourceRepo(dbx)
	now := time.Now()

	mock.ExpectQuery("SELECT id, owner_user_id, category_id, title, description, location, price_per_hour, is_active, created_at FROM resources WHERE owner_user_id = \\? ORDER BY id DESC").
		WithArgs(uint64(9)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "owner_user_id", "category_id", "title", "description", "location", "price_per_hour", "is_active", "created_at",
		}).AddRow(uint64(1), uint64(9), uint64(1), "Mine", nil, nil, 0, true, now))

	_, err := r.ListByOwner(context.Background(), 9)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
