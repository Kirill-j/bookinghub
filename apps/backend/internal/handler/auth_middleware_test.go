package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/service"
)

func TestAuthMiddleware_NoHeader(t *testing.T) {
	auth := service.NewAuthService("secret", 15)

	h := AuthMiddleware(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != 401 {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestAuthMiddleware_BadToken(t *testing.T) {
	auth := service.NewAuthService("secret", 15)

	h := AuthMiddleware(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer BADTOKEN")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != 401 {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestAuthMiddleware_OK_PutsClaimsToContext(t *testing.T) {
	auth := service.NewAuthService("secret", 15)

	tok, err := auth.CreateAccessToken(123, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("CreateAccessToken: %v", err)
	}

	called := false
	h := AuthMiddleware(auth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if GetUserID(r) != 123 {
			t.Fatalf("expected uid=123 got %d", GetUserID(r))
		}
		if GetRole(r) != domain.RoleAdmin {
			t.Fatalf("expected role ADMIN got %s", GetRole(r))
		}
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler called")
	}
	if rr.Code != 200 {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestRequireRoles_Forbidden(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	// сделаем запрос без роли в контексте -> 401
	h := RequireRoles(domain.RoleAdmin)(next)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != 401 {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}
