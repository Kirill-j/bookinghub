package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newRepoMock2(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock, func() { _ = db.Close() }
}

func TestCategoryRepo_List_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock2(t)
	defer cleanup()

	r := NewCategoryRepo(dbx)
	now := time.Now()

	mock.ExpectQuery("SELECT id, name, created_at FROM resource_categories ORDER BY name ASC").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(uint64(1), "A", now))

	_, err := r.List(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCategoryRepo_Create_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock2(t)
	defer cleanup()

	r := NewCategoryRepo(dbx)

	mock.ExpectExec("INSERT INTO resource_categories \\(name\\) VALUES \\(\\?\\)").
		WithArgs("X").
		WillReturnResult(sqlmock.NewResult(7, 1))

	id, err := r.Create(context.Background(), "X")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if id != 7 {
		t.Fatalf("expected 7 got %d", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCategoryRepo_Update_OK(t *testing.T) {
	dbx, mock, cleanup := newRepoMock2(t)
	defer cleanup()

	r := NewCategoryRepo(dbx)

	mock.ExpectExec("UPDATE resource_categories SET name = \\? WHERE id = \\?").
		WithArgs("Y", uint64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.Update(context.Background(), 2, "Y"); err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCategoryRepo_Delete_Error(t *testing.T) {
	dbx, mock, cleanup := newRepoMock2(t)
	defer cleanup()

	r := NewCategoryRepo(dbx)

	mock.ExpectExec("DELETE FROM resource_categories WHERE id = \\?").
		WithArgs(uint64(2)).
		WillReturnError(errors.New("fk"))

	if err := r.Delete(context.Background(), 2); err == nil {
		t.Fatalf("expected error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
