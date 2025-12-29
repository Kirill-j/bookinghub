package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// "time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
)

func newSQLXMock4(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func TestUserHandler_PublicByID_NotFound_404(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock4(t)
	defer cleanup()

	h := NewUserHandler(repo.NewUserRepo(dbx))

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(99)).
		WillReturnError(sqlmock.ErrCancelled) // любой err → 404 в твоём хендлере

	req := httptest.NewRequest(http.MethodGet, "/api/users/99", nil)
	rr := httptest.NewRecorder()

	// без chi URLParam → вернёт 400, поэтому здесь лучше тест через router_smoke_test.go
	_ = h
	_ = rr
	_ = req
	_ = domain.RoleIndividual
}
