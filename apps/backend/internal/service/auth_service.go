package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"bookinghub-backend/internal/domain"
)

type AuthService struct {
	jwtSecret []byte
	accessTTL time.Duration
}

func NewAuthService(jwtSecret string, accessTTLMinutes int) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
		accessTTL: time.Duration(accessTTLMinutes) * time.Minute,
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (s *AuthService) CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type AccessClaims struct {
	UserID uint64          `json:"userId"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) CreateAccessToken(userID uint64, role domain.UserRole) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			Subject:   strconv.FormatUint(userID, 10),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*AccessClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
