package repo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"bookinghub-backend/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepo_GetRoleByID(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewUserRepo(db)

	q := regexp.QuoteMeta(`
		SELECT role
		FROM users
		WHERE id = ?
		LIMIT 1
	`)
	rows := sqlmock.NewRows([]string{"role"}).AddRow(string(domain.RoleAdmin))
	mock.ExpectQuery(q).WithArgs(uint64(1)).WillReturnRows(rows)

	role, err := r.GetRoleByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRoleByID err: %v", err)
	}
	if role != domain.RoleAdmin {
		t.Fatalf("expected ADMIN got %s", role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestUserRepo_GetByEmail_OK(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewUserRepo(db)
	now := time.Date(2025, 12, 29, 12, 0, 0, 0, time.UTC)

	q := regexp.QuoteMeta(`
		SELECT id, email, name, role, password_hash, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`)

	rows := sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
		AddRow(uint64(10), "a@a.ru", "A", string(domain.RoleIndividual), "HASH", now)

	mock.ExpectQuery(q).WithArgs("a@a.ru").WillReturnRows(rows)

	u, err := r.GetByEmail(context.Background(), "a@a.ru")
	if err != nil {
		t.Fatalf("GetByEmail err: %v", err)
	}
	if u.ID != 10 || u.Email != "a@a.ru" {
		t.Fatalf("unexpected user: %+v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestUserRepo_Create(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewUserRepo(db)

	q := regexp.QuoteMeta(`
		INSERT INTO users (email, name, role, password_hash)
		VALUES (?, ?, ?, ?)
	`)
	mock.ExpectExec(q).
		WithArgs("x@x.ru", "X", domain.RoleCompany, "HASH").
		WillReturnResult(sqlmock.NewResult(77, 1))

	id, err := r.Create(context.Background(), "x@x.ru", "X", domain.RoleCompany, "HASH")
	if err != nil {
		t.Fatalf("Create err: %v", err)
	}
	if id != 77 {
		t.Fatalf("expected 77 got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestUserRepo_DeleteAccount(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewUserRepo(db)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM bookings WHERE user_id = ?`)).
		WithArgs(uint64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE b FROM bookings b
		JOIN resources r ON r.id = b.resource_id
		WHERE r.owner_user_id = ?
	`)).
		WithArgs(uint64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM resources WHERE owner_user_id = ?`)).
		WithArgs(uint64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = ?`)).
		WithArgs(uint64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	if err := r.DeleteAccount(context.Background(), 5); err != nil {
		t.Fatalf("DeleteAccount err: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	r := NewUserRepo(db)

	q := regexp.QuoteMeta(`
		SELECT id, email, name, role, password_hash, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`)
	mock.ExpectQuery(q).WithArgs("no@no.ru").WillReturnError(sql.ErrNoRows)

	u, err := r.GetByEmail(context.Background(), "no@no.ru")
	if err == nil {
		t.Fatalf("expected error")
	}
	if u != nil {
		t.Fatalf("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
