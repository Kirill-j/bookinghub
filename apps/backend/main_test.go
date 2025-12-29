package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestGetEnv(t *testing.T) {
	_ = os.Unsetenv("X_TEST_ENV")
	if got := getEnv("X_TEST_ENV", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}

	_ = os.Setenv("X_TEST_ENV", "abc")
	if got := getEnv("X_TEST_ENV", "fallback"); got != "abc" {
		t.Fatalf("expected abc, got %q", got)
	}
}

func TestHandleDBPing_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectPing()

	app := &App{DB: sqlx.NewDb(db, "sqlmock")}

	req := httptest.NewRequest(http.MethodGet, "/db/ping", nil)
	rr := httptest.NewRecorder()

	app.handleDBPing(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}
