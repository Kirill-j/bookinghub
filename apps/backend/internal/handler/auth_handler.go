package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"
)

type AuthHandler struct {
	users *repo.UserRepo
	auth  *service.AuthService
}

func NewAuthHandler(users *repo.UserRepo, auth *service.AuthService) *AuthHandler {
	return &AuthHandler{users: users, auth: auth}
}

type registerReq struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	AccountType string `json:"accountType"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, "Введите корректный email", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Имя обязательно", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Пароль должен быть не короче 6 символов", http.StatusBadRequest)
		return
	}

	// проверим, что email не занят
	existing, err := h.users.GetByEmail(r.Context(), req.Email)
	if err == nil && existing != nil {
		http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
		return
	}
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := h.auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Не удалось обработать пароль", http.StatusInternalServerError)
		return
	}

	role := domain.RoleIndividual
	switch strings.ToUpper(strings.TrimSpace(req.AccountType)) {
	case "INDIVIDUAL", "":
		role = domain.RoleIndividual
	case "COMPANY":
		role = domain.RoleCompany
	default:
		http.Error(w, "Некорректный тип аккаунта (accountType)", http.StatusBadRequest)
		return
	}

	id, err := h.users.Create(r.Context(), req.Email, req.Name, role, hash)
	if err != nil {
		http.Error(w, "Не удалось создать пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.auth.CreateAccessToken(id, role)
	if err != nil {
		http.Error(w, "Не удалось создать токен", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"accessToken": token,
		"user": map[string]any{
			"id":    id,
			"email": req.Email,
			"name":  req.Name,
			"role":  role,
		},
	})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	u, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Если это сид-пользователь с TEMP — зададим ему пароль один раз:
	// admin/manager/user -> пароль "123456"
	if u.PasswordHash == "TEMP" {
		hash, _ := h.auth.HashPassword("123456")
		_ = h.users.UpdatePasswordHash(r.Context(), u.Email, hash)
		u.PasswordHash = hash
	}

	if err := h.auth.CheckPassword(u.PasswordHash, req.Password); err != nil {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	token, err := h.auth.CreateAccessToken(u.ID, u.Role)
	if err != nil {
		http.Error(w, "Не удалось создать токен", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken": token,
		"user": map[string]any{
			"id":    u.ID,
			"email": u.Email,
			"name":  u.Name,
			"role":  u.Role,
		},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	u, err := h.users.GetByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"role":  u.Role,
	})
}

type updateMeReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	var req updateMeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	name := strings.TrimSpace(req.Name)

	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, "Введите корректный email", http.StatusBadRequest)
		return
	}
	if name == "" {
		http.Error(w, "Имя обязательно", http.StatusBadRequest)
		return
	}

	// текущий пользователь
	u, err := h.users.GetByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// если меняем email — проверим уникальность
	if email != u.Email {
		existing, err := h.users.GetByEmail(r.Context(), email)
		if err == nil && existing != nil {
			http.Error(w, "Email уже занят", http.StatusConflict)
			return
		}
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Ошибка базы данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := h.users.UpdateProfile(r.Context(), uid, email, name); err != nil {
		http.Error(w, "Не удалось обновить профиль: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// отдадим обновлённого пользователя
	u2, _ := h.users.GetByID(r.Context(), uid)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":    u2.ID,
		"email": u2.Email,
		"name":  u2.Name,
		"role":  u2.Role,
	})
}

type changePasswordReq struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	var req changePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 6 {
		http.Error(w, "Новый пароль должен быть минимум 6 символов", http.StatusBadRequest)
		return
	}

	u, err := h.users.GetByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	if err := h.auth.CheckPassword(u.PasswordHash, req.CurrentPassword); err != nil {
		http.Error(w, "Текущий пароль неверный", http.StatusUnauthorized)
		return
	}

	hash, err := h.auth.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Не удалось обработать пароль", http.StatusInternalServerError)
		return
	}

	if err := h.users.UpdatePasswordHashByID(r.Context(), uid, hash); err != nil {
		http.Error(w, "Не удалось обновить пароль: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *AuthHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	if err := h.users.DeleteAccount(r.Context(), uid); err != nil {
		http.Error(w, "Не удалось удалить аккаунт: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
