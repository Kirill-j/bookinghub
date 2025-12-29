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
	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/repo"
)

func newSQLXMock2(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func withUIDRes(ctx context.Context, uid uint64) context.Context {
	return context.WithValue(ctx, ctxUserID, uid)
}

func TestResourceHandler_List_OK(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock2(t)
	defer cleanup()

	h := NewResourceHandler(repo.NewResourceRepo(dbx))

	now := time.Now()
	mock.ExpectQuery("SELECT id, owner_user_id, category_id, title, description, location, price_per_hour, is_active, created_at FROM resources ORDER BY id DESC").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "owner_user_id", "category_id", "title", "description", "location", "price_per_hour", "is_active", "created_at",
		}).AddRow(uint64(1), uint64(2), uint64(3), "Title", nil, nil, 100, true, now))

	req := httptest.NewRequest(http.MethodGet, "/api/resources", nil)
	rr := httptest.NewRecorder()

	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestResourceHandler_My_Unauthorized_401(t *testing.T) {
	dbx, _, cleanup := newSQLXMock2(t)
	defer cleanup()

	h := NewResourceHandler(repo.NewResourceRepo(dbx))

	req := httptest.NewRequest(http.MethodGet, "/api/resources/my", nil)
	rr := httptest.NewRecorder()

	h.My(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestResourceHandler_My_OK(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock2(t)
	defer cleanup()

	h := NewResourceHandler(repo.NewResourceRepo(dbx))

	now := time.Now()
	mock.ExpectQuery("SELECT id, owner_user_id, category_id, title, description, location, price_per_hour, is_active, created_at FROM resources WHERE owner_user_id = \\? ORDER BY id DESC").
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "owner_user_id", "category_id", "title", "description", "location", "price_per_hour", "is_active", "created_at",
		}).AddRow(uint64(10), uint64(5), uint64(1), "Mine", nil, nil, 0, true, now))

	req := httptest.NewRequest(http.MethodGet, "/api/resources/my", nil)
	req = req.WithContext(withUIDRes(req.Context(), 5))
	rr := httptest.NewRecorder()

	h.My(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestResourceHandler_Create_Validation_400(t *testing.T) {
	dbx, _, cleanup := newSQLXMock2(t)
	defer cleanup()

	h := NewResourceHandler(repo.NewResourceRepo(dbx))

	body := map[string]any{"categoryId": 0, "title": ""}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/resources", bytes.NewReader(b))
	req = req.WithContext(withUIDRes(req.Context(), 5))
	rr := httptest.NewRecorder()

	h.Create(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestResourceHandler_Create_OK_201(t *testing.T) {
	dbx, mock, cleanup := newSQLXMock2(t)
	defer cleanup()

	h := NewResourceHandler(repo.NewResourceRepo(dbx))

	mock.ExpectExec("INSERT INTO resources \\(owner_user_id, category_id, title, description, location, price_per_hour\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(uint64(7), uint64(2), "Hello", nil, nil, 100).
		WillReturnResult(sqlmock.NewResult(55, 1))

	body := map[string]any{
		"categoryId":   2,
		"title":        "Hello",
		"description":  nil,
		"location":     nil,
		"pricePerHour": 100,
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/resources", bytes.NewReader(b))
	req = req.WithContext(withUIDRes(req.Context(), 7))
	rr := httptest.NewRecorder()

	h.Create(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
