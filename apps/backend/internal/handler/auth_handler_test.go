package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"
)

func newSQLXMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cleanup := func() {
		_ = db.Close()
	}
	return sqlxDB, mock, cleanup
}

func withUser(ctx context.Context, uid uint64, role domain.UserRole) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, uid)
	ctx = context.WithValue(ctx, ctxRole, role)
	return ctx
}

func TestAuthHandler_Register_BadJSON(t *testing.T) {
	dbx, _, cleanup := newSQLXMock(t)
	defer cleanup()

	h := NewAuthHandler(repo.NewUserRepo(dbx), service.NewAuthService("dev", 15))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("{bad"))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Register_ExistingEmail_409(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	auth := service.NewAuthService("dev", 15)
	h := NewAuthHandler(users, auth)

	created := time.Now()

	// GetByEmail -> found
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\? LIMIT 1").
		WithArgs("x@test.local").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
				AddRow(uint64(1), "x@test.local", "X", string(domain.RoleIndividual), "HASH", created),
		)

	body := map[string]any{
		"email":       "x@test.local",
		"name":        "X",
		"password":    "123456",
		"accountType": "INDIVIDUAL",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_Register_OK_201(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	auth := service.NewAuthService("dev", 15)
	h := NewAuthHandler(users, auth)

	// GetByEmail -> sql.ErrNoRows
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\? LIMIT 1").
		WithArgs("new@test.local").
		WillReturnError(sql.ErrNoRows)

	// Create -> insert id
	mock.ExpectExec("INSERT INTO users \\(email, name, role, password_hash\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
		WithArgs("new@test.local", "New", string(domain.RoleCompany), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(7, 1))

	body := map[string]any{
		"email":       "new@test.local",
		"name":        "New",
		"password":    "123456",
		"accountType": "COMPANY",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_Login_OK_200(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	auth := service.NewAuthService("dev", 15)
	h := NewAuthHandler(users, auth)

	hash, _ := auth.HashPassword("123456")
	created := time.Now()

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\? LIMIT 1").
		WithArgs("a@test.local").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(10), "a@test.local", "A", string(domain.RoleIndividual), hash, created))

	body := map[string]any{"email": "a@test.local", "password": "123456"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_Login_TEMPUser_SetsDefaultPassword(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	auth := service.NewAuthService("dev", 15)
	h := NewAuthHandler(users, auth)

	created := time.Now()

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\? LIMIT 1").
		WithArgs("temp@test.local").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(11), "temp@test.local", "Temp", string(domain.RoleIndividual), "TEMP", created))

	// UpdatePasswordHash(email, hash)
	mock.ExpectExec("UPDATE users SET password_hash = \\? WHERE email = \\?").
		WithArgs(sqlmock.AnyArg(), "temp@test.local").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := map[string]any{"email": "temp@test.local", "password": "123456"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_Me_Unauthorized_401(t *testing.T) {
	dbx, _, cleanup := newSQLXMock(t)
	defer cleanup()

	h := NewAuthHandler(repo.NewUserRepo(dbx), service.NewAuthService("dev", 15))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rr := httptest.NewRecorder()

	h.Me(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Me_OK_200(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	h := NewAuthHandler(users, service.NewAuthService("dev", 15))

	created := time.Now()

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(5), "me@test.local", "Me", string(domain.RoleCompany), "HASH", created))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req = req.WithContext(withUser(req.Context(), 5, domain.RoleCompany))
	rr := httptest.NewRecorder()

	h.Me(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_UpdateMe_OK_200_EmailChanged(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	h := NewAuthHandler(users, service.NewAuthService("dev", 15))

	created := time.Now()

	// current user
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(5), "old@test.local", "Old", string(domain.RoleIndividual), "HASH", created))

	// uniqueness: GetByEmail(new) -> no rows
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\? LIMIT 1").
		WithArgs("new@test.local").
		WillReturnError(sql.ErrNoRows)

	// UpdateProfile
	mock.ExpectExec("UPDATE users SET email = \\?, name = \\? WHERE id = \\?").
		WithArgs("new@test.local", "NewName", uint64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// GetByID again
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(5), "new@test.local", "NewName", string(domain.RoleIndividual), "HASH", created))

	body := map[string]any{"email": "new@test.local", "name": "NewName"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/auth/me", bytes.NewReader(b))
	req = req.WithContext(withUser(req.Context(), 5, domain.RoleIndividual))
	rr := httptest.NewRecorder()

	h.UpdateMe(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_ChangePassword_OK_200(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	auth := service.NewAuthService("dev", 15)
	h := NewAuthHandler(users, auth)

	oldHash, _ := auth.HashPassword("oldpass")
	created := time.Now()

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(9), "p@test.local", "P", string(domain.RoleIndividual), oldHash, created))

	mock.ExpectExec("UPDATE users SET password_hash = \\? WHERE id = \\?").
		WithArgs(sqlmock.AnyArg(), uint64(9)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := map[string]any{"currentPassword": "oldpass", "newPassword": "newpass1"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/password", bytes.NewReader(b))
	req = req.WithContext(withUser(req.Context(), 9, domain.RoleIndividual))
	rr := httptest.NewRecorder()

	h.ChangePassword(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestAuthHandler_DeleteMe_OK_200(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock(t)
	defer cleanup()

	users := repo.NewUserRepo(dbx)
	h := NewAuthHandler(users, service.NewAuthService("dev", 15))

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM bookings WHERE user_id = \\?").WithArgs(uint64(3)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE b FROM bookings b JOIN resources r ON r.id = b.resource_id WHERE r.owner_user_id = \\?").
		WithArgs(uint64(3)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM resources WHERE owner_user_id = \\?").WithArgs(uint64(3)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM users WHERE id = \\?").WithArgs(uint64(3)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/api/auth/me", nil)
	req = req.WithContext(withUser(req.Context(), 3, domain.RoleIndividual))
	rr := httptest.NewRecorder()

	h.DeleteMe(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
