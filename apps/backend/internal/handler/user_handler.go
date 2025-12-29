package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"bookinghub-backend/internal/repo"
)

type UserHandler struct {
	users *repo.UserRepo
}

func NewUserHandler(users *repo.UserRepo) *UserHandler {
	return &UserHandler{users: users}
}

// GET /api/users/{id}
func (h *UserHandler) PublicByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	u, err := h.users.GetByID(r.Context(), uint64(id64))
	if err != nil || u == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// публичные поля (можно скрыть email если хочешь)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":        u.ID,
		"name":      u.Name,
		"role":      u.Role,
		"email":     u.Email, // если не хочешь светить — убери
		"createdAt": u.CreatedAt,
	})
}
