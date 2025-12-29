package repo

import (
	"context"
	"testing"
	"time"

	"bookinghub-backend/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newSQLXMockRepo2(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func TestUserRepo_GetByID_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo2(t)
	defer closeFn()

	r := NewUserRepo(db)

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\?").
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(5), "a@b.c", "Alex", "INDIVIDUAL", "hash", time.Now()))

	u, err := r.GetByID(context.Background(), 5)
	if err != nil || u == nil {
		t.Fatalf("expected user, got err=%v u=%v", err, u)
	}
	if u.Email != "a@b.c" || u.Role != domain.RoleIndividual {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestUserRepo_UpdateProfile_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo2(t)
	defer closeFn()

	r := NewUserRepo(db)

	mock.ExpectExec("UPDATE users\\s+SET email = \\?, name = \\?\\s+WHERE id = \\?").
		WithArgs("new@b.c", "NewName", uint64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.UpdateProfile(context.Background(), 7, "new@b.c", "NewName"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestUserRepo_UpdatePasswordHash_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo2(t)
	defer closeFn()

	r := NewUserRepo(db)

	mock.ExpectExec("UPDATE users SET password_hash = \\?\\s+WHERE email = \\?").
		WithArgs("h2", "x@y.z").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.UpdatePasswordHash(context.Background(), "x@y.z", "h2"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestUserRepo_UpdatePasswordHashByID_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMockRepo2(t)
	defer closeFn()

	r := NewUserRepo(db)

	mock.ExpectExec("UPDATE users\\s+SET password_hash = \\?\\s+WHERE id = \\?").
		WithArgs("h3", uint64(9)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := r.UpdatePasswordHashByID(context.Background(), 9, "h3"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
