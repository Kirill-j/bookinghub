package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"
)

func newMockHandlerDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, mock, func() { _ = db.Close() }
}

func withUID(req *http.Request, uid uint64) *http.Request {
	ctx := context.WithValue(req.Context(), ctxUserID, uid)
	return req.WithContext(ctx)
}

func TestBookingHandler_Create_BadJSON(t *testing.T) {
	db, mock, cleanup := newMockHandlerDB(t)
	defer cleanup()

	bookingRepo := repo.NewBookingRepo(db)
	userRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bookingRepo)
	h := NewBookingHandler(bookingRepo, userRepo, svc)

	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewBufferString("{bad"))
	req = withUID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)
	if rr.Code != 400 {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
	_ = mock.ExpectationsWereMet()
}

func TestBookingHandler_Create_Conflict(t *testing.T) {
	db, mock, cleanup := newMockHandlerDB(t)
	defer cleanup()

	bookingRepo := repo.NewBookingRepo(db)
	userRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bookingRepo)
	h := NewBookingHandler(bookingRepo, userRepo, svc)

	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	end := start.Add(time.Hour)

	body, _ := json.Marshal(map[string]any{
		"resourceId": 99,
		"startAt":    start.Format(time.RFC3339),
		"endAt":      end.Format(time.RFC3339),
	})

	// service.Create -> repo.HasConflict -> COUNT(*) = 1
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM bookings
		WHERE resource_id = ?
		  AND status IN ('PENDING','APPROVED')
		  AND (? < end_at) AND (? > start_at)
	`)).
		WithArgs(uint64(99), start, end).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))

	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body))
	req = withUID(req, 7)
	rr := httptest.NewRecorder()

	h.Create(rr, req)
	if rr.Code != 409 {
		t.Fatalf("expected 409 got %d body=%s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingHandler_Create_OK(t *testing.T) {
	db, mock, cleanup := newMockHandlerDB(t)
	defer cleanup()

	bookingRepo := repo.NewBookingRepo(db)
	userRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bookingRepo)
	h := NewBookingHandler(bookingRepo, userRepo, svc)

	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	end := start.Add(time.Hour)

	body, _ := json.Marshal(map[string]any{
		"resourceId": 99,
		"startAt":    start.Format(time.RFC3339),
		"endAt":      end.Format(time.RFC3339),
	})

	// no conflict
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM bookings
		WHERE resource_id = ?
		  AND status IN ('PENDING','APPROVED')
		  AND (? < end_at) AND (? > start_at)
	`)).
		WithArgs(uint64(99), start, end).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))

	// insert booking
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO bookings (resource_id, user_id, start_at, end_at, status)
		VALUES (?, ?, ?, ?, 'PENDING')
	`)).
		WithArgs(uint64(99), uint64(7), start, end).
		WillReturnResult(sqlmock.NewResult(555, 1))

	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body))
	req = withUID(req, 7)
	rr := httptest.NewRecorder()

	h.Create(rr, req)
	if rr.Code != 201 {
		t.Fatalf("expected 201 got %d body=%s", rr.Code, rr.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestBookingHandler_Pending_Admin(t *testing.T) {
	db, mock, cleanup := newMockHandlerDB(t)
	defer cleanup()

	bookingRepo := repo.NewBookingRepo(db)
	userRepo := repo.NewUserRepo(db)
	svc := service.NewBookingService(bookingRepo)
	h := NewBookingHandler(bookingRepo, userRepo, svc)

	// GetRoleByID -> ADMIN
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT role
		FROM users
		WHERE id = ?
		LIMIT 1
	`)).
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(string(domain.RoleAdmin)))

	// ListPending
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, resource_id, user_id, start_at, end_at, status, manager_comment, created_at, updated_at
		FROM bookings
		WHERE status = 'PENDING'
		ORDER BY start_at ASC
	`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "resource_id", "user_id", "start_at", "end_at", "status", "manager_comment", "created_at", "updated_at",
		}).AddRow(uint64(1), uint64(2), uint64(3), now, now.Add(time.Hour), "PENDING", nil, now, nil))

	req := httptest.NewRequest("GET", "/api/bookings/pending", nil)
	req = withUID(req, 1)
	rr := httptest.NewRecorder()

	h.Pending(rr, req)
	if rr.Code != 200 {
		t.Fatalf("expected 200 got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
