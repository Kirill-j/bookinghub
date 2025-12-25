package handler

import (
	"context"
	"net/http"
	"strings"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/service"
)

type ctxKey string

const (
	ctxUserID ctxKey = "userId"
	ctxRole   ctxKey = "role"
)

func AuthMiddleware(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(h, "Bearer ")

			claims, err := auth.ParseAccessToken(tokenStr)
			if err != nil {
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) uint64 {
	v := r.Context().Value(ctxUserID)
	if v == nil {
		return 0
	}
	return v.(uint64)
}

func GetRole(r *http.Request) domain.UserRole {
	v := r.Context().Value(ctxRole)
	if v == nil {
		return ""
	}
	return v.(domain.UserRole)
}

// RequireRoles — middleware, которое пускает только нужные роли
func RequireRoles(allowed ...domain.UserRole) func(http.Handler) http.Handler {
	allowedSet := map[domain.UserRole]struct{}{}
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetRole(r)
			if role == "" {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}
			if _, ok := allowedSet[role]; !ok {
				http.Error(w, "Недостаточно прав", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
