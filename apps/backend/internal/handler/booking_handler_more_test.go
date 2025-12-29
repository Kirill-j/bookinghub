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
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func newSQLXMock2Res(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() { _ = db.Close() }
}

func withUIDBH(ctx context.Context, uid uint64) context.Context {
	return context.WithValue(ctx, ctxUserID, uid)
}

func TestBookingHandler_My_Unauthorized(t *testing.T) {
	db, _, closeFn := newSQLXMock2Res(t)
	defer closeFn()

	bRepo := repo.NewBookingRepo(db)
	uRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bRepo)
	h := NewBookingHandler(bRepo, uRepo, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/bookings/my", nil)
	rr := httptest.NewRecorder()

	h.My(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBookingHandler_My_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMock2Res(t)
	defer closeFn()

	bRepo := repo.NewBookingRepo(db)
	uRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bRepo)
	h := NewBookingHandler(bRepo, uRepo, svc)

	mock.ExpectQuery("SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at").
		WithArgs(uint64(10)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "resource_id", "user_id", "start_at", "end_at", "status", "manager_comment", "created_at", "updated_at",
		}).AddRow(
			uint64(1), uint64(2), uint64(10),
			time.Now().Add(2*time.Hour), time.Now().Add(3*time.Hour),
			"PENDING", nil, time.Now(), nil,
		))

	req := httptest.NewRequest(http.MethodGet, "/api/bookings/my", nil)
	req = req.WithContext(withUIDBH(req.Context(), 10))
	rr := httptest.NewRecorder()

	h.My(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookingHandler_UpdateStatus_Forbidden_NotOwner(t *testing.T) {
	db, mock, closeFn := newSQLXMock2Res(t)
	defer closeFn()

	bRepo := repo.NewBookingRepo(db)
	uRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bRepo)
	h := NewBookingHandler(bRepo, uRepo, svc)

	// owner of booking -> 999, current user -> 10
	mock.ExpectQuery("SELECT r.owner_user_id").
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"owner_user_id"}).AddRow(uint64(999)))

	// role current user -> INDIVIDUAL
	mock.ExpectQuery("SELECT role FROM users").
		WithArgs(uint64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("INDIVIDUAL"))

	body, _ := json.Marshal(map[string]any{
		"status":         "APPROVED",
		"managerComment": nil,
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/bookings/1/status", bytes.NewReader(body))
	req = req.WithContext(withUIDBH(req.Context(), 10))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.UpdateStatus(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBookingHandler_UpdateStatus_OK_Owner(t *testing.T) {
	db, mock, closeFn := newSQLXMock2Res(t)
	defer closeFn()

	bRepo := repo.NewBookingRepo(db)
	uRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bRepo)
	h := NewBookingHandler(bRepo, uRepo, svc)

	// owner is current user
	mock.ExpectQuery("SELECT r.owner_user_id").
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"owner_user_id"}).AddRow(uint64(10)))

	// role current user
	mock.ExpectQuery("SELECT role FROM users").
		WithArgs(uint64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("INDIVIDUAL"))

	// booking exists and pending
	mock.ExpectQuery("SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at").
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "resource_id", "user_id", "start_at", "end_at", "status", "manager_comment", "created_at", "updated_at",
		}).AddRow(
			uint64(7), uint64(2), uint64(55),
			time.Now().Add(2*time.Hour), time.Now().Add(3*time.Hour),
			"PENDING", nil, time.Now(), nil,
		))

	// update status
	mock.ExpectExec("UPDATE bookings\\s+SET status = \\?, manager_comment = \\?\\s+WHERE id = \\?").
		WithArgs("APPROVED", sqlmock.AnyArg(), uint64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]any{
		"status":         "APPROVED",
		"managerComment": "ok",
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/bookings/7/status", bytes.NewReader(body))
	req = req.WithContext(withUIDBH(req.Context(), 10))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "7")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.UpdateStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookingHandler_Cancel_OK(t *testing.T) {
	db, mock, closeFn := newSQLXMock2Res(t)
	defer closeFn()

	bRepo := repo.NewBookingRepo(db)
	uRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bRepo)
	h := NewBookingHandler(bRepo, uRepo, svc)

	start := time.Now().Add(5 * time.Hour)

	// booking exists, belongs to user, status pending
	mock.ExpectQuery("SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at").
		WithArgs(uint64(3)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "resource_id", "user_id", "start_at", "end_at", "status", "manager_comment", "created_at", "updated_at",
		}).AddRow(
			uint64(3), uint64(2), uint64(10),
			start, start.Add(time.Hour),
			"PENDING", nil, time.Now(), nil,
		))

	mock.ExpectExec("UPDATE bookings\\s+SET status = 'CANCELED'\\s+WHERE id = \\?").
		WithArgs(uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodPost, "/api/bookings/3/cancel", nil)
	req = req.WithContext(withUIDBH(req.Context(), 10))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "3")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h.Cancel(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
