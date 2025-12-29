package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
)

func newSQLXMock5(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func TestRouter_Category_Update_OK(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock5(t)
	defer cleanup()

	catH := NewCategoryHandler(repo.NewCategoryRepo(dbx))

	mock.ExpectExec("UPDATE resource_categories SET name = \\? WHERE id = \\?").
		WithArgs("X", uint64(12)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := chi.NewRouter()
	r.Patch("/api/categories/{id}", catH.Update)

	body := map[string]any{"name": "X"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/categories/12", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRouter_Category_Delete_OK(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock5(t)
	defer cleanup()

	catH := NewCategoryHandler(repo.NewCategoryRepo(dbx))

	mock.ExpectExec("DELETE FROM resource_categories WHERE id = \\?").
		WithArgs(uint64(12)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := chi.NewRouter()
	r.Delete("/api/categories/{id}", catH.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/12", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRouter_User_PublicByID_OK(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock5(t)
	defer cleanup()

	userH := NewUserHandler(repo.NewUserRepo(dbx))

	now := time.Now()
	mock.ExpectQuery("SELECT id, email, name, role, password_hash, created_at FROM users WHERE id = \\? LIMIT 1").
		WithArgs(uint64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "role", "password_hash", "created_at"}).
			AddRow(uint64(2), "u@test.local", "User", string(domain.RoleIndividual), "HASH", now))

	r := chi.NewRouter()
	r.Get("/api/users/{id}", userH.PublicByID)

	req := httptest.NewRequest(http.MethodGet, "/api/users/2", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRouter_Resource_Create_Unauthorized_401(t *testing.T) {
	dbx, _, cleanup := newSQLXMock5(t)
	defer cleanup()

	resH := NewResourceHandler(repo.NewResourceRepo(dbx))

	r := chi.NewRouter()
	r.Post("/api/resources", resH.Create)

	body := map[string]any{"categoryId": 1, "title": "X", "pricePerHour": 0}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/resources", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRouter_Resource_Create_OK_201_WithContextUser(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock5(t)
	defer cleanup()

	resH := NewResourceHandler(repo.NewResourceRepo(dbx))

	mock.ExpectExec("INSERT INTO resources \\(owner_user_id, category_id, title, description, location, price_per_hour\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(uint64(9), uint64(1), "X", nil, nil, 0).
		WillReturnResult(sqlmock.NewResult(101, 1))

	r := chi.NewRouter()
	// имитируем auth: просто кладём userId в context
	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), ctxUserID, uint64(9))
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}).Post("/api/resources", resH.Create)

	body := map[string]any{"categoryId": 1, "title": "X", "pricePerHour": 0}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/resources", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
