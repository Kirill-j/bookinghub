package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newSQLXMockAuth(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func withUIDAuth(ctx context.Context, uid uint64) context.Context {
	return context.WithValue(ctx, ctxUserID, uid)
}

func TestAuthHandler_Register_BadJSON_New(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("{bad"))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	body, _ := json.Marshal(map[string]any{
		"email":       "no-at",
		"name":        "A",
		"password":    "123456",
		"accountType": "INDIVIDUAL",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	body, _ := json.Marshal(map[string]any{
		"email":    "a@b.c",
		"name":     "A",
		"password": "123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Register_ConflictEmail(t *testing.T) {
	db, mock, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	// users.GetByEmail -> returns existing row
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\?").
		WithArgs("a@b.c").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(1), "a@b.c", "Alex", "INDIVIDUAL", "hash", time.Now()))

	body, _ := json.Marshal(map[string]any{
		"email":    "a@b.c",
		"name":     "A",
		"password": "123456",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Login_BadJSON(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("{bad"))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Login_EmptyFields(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	body, _ := json.Marshal(map[string]any{"email": "", "password": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	db, mock, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	hash, _ := auth.HashPassword("correct123")

	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE email = \\?").
		WithArgs("a@b.c").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(1), "a@b.c", "Alex", "INDIVIDUAL", hash, time.Now()))

	body, _ := json.Marshal(map[string]any{"email": "a@b.c", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_Me_Unauthorized(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rr := httptest.NewRecorder()

	h.Me(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_UpdateMe_InvalidJSON(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	req := httptest.NewRequest(http.MethodPatch, "/api/auth/me", bytes.NewBufferString("{bad"))
	req = req.WithContext(withUIDAuth(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.UpdateMe(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_UpdateMe_InvalidEmail(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	body, _ := json.Marshal(map[string]any{"email": "bad", "name": "A"})
	req := httptest.NewRequest(http.MethodPatch, "/api/auth/me", bytes.NewReader(body))
	req = req.WithContext(withUIDAuth(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.UpdateMe(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_ChangePassword_TooShort(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	body, _ := json.Marshal(map[string]any{
		"currentPassword": "123456",
		"newPassword":     "123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewReader(body))
	req = req.WithContext(withUIDAuth(req.Context(), 1))
	rr := httptest.NewRecorder()

	h.ChangePassword(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthHandler_DeleteMe_Unauthorized(t *testing.T) {
	db, _, closeFn := newSQLXMockAuth(t)
	defer closeFn()

	users := repo.NewUserRepo(db)
	auth := service.NewAuthService("secret", 60)
	h := NewAuthHandler(users, auth)

	req := httptest.NewRequest(http.MethodDelete, "/api/auth/me", nil)
	rr := httptest.NewRecorder()

	h.DeleteMe(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}
