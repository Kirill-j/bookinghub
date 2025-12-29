package service

import (
	"testing"

	"bookinghub-backend/internal/domain"
)

func TestAuthService_HashAndCheckPassword_OK(t *testing.T) {
	s := NewAuthService("test-secret", 15)

	hash, err := s.HashPassword("123456")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}
	if hash == "" {
		t.Fatalf("hash empty")
	}

	if err := s.CheckPassword(hash, "123456"); err != nil {
		t.Fatalf("check should pass, got: %v", err)
	}
	if err := s.CheckPassword(hash, "wrong"); err == nil {
		t.Fatalf("check should fail for wrong password")
	}
}

func TestAuthService_Token_Roundtrip_OK(t *testing.T) {
	s := NewAuthService("test-secret", 15)

	token, err := s.CreateAccessToken(42, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("create token err: %v", err)
	}

	claims, err := s.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("parse token err: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected userId=42, got %d", claims.UserID)
	}
	if claims.Role != domain.RoleAdmin {
		t.Fatalf("expected role=%s, got %s", domain.RoleAdmin, claims.Role)
	}
	if claims.Subject != "42" {
		t.Fatalf("expected subject=42, got %s", claims.Subject)
	}
}

func TestAuthService_ParseToken_Invalid(t *testing.T) {
	s := NewAuthService("test-secret", 15)

	_, err := s.ParseAccessToken("not-a-token")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAuthService_ParseToken_WrongSecret(t *testing.T) {
	s1 := NewAuthService("secret-1", 15)
	s2 := NewAuthService("secret-2", 15)

	token, err := s1.CreateAccessToken(1, domain.RoleIndividual)
	if err != nil {
		t.Fatalf("create token err: %v", err)
	}

	_, err = s2.ParseAccessToken(token)
	if err == nil {
		t.Fatalf("expected error due to wrong secret")
	}
}
