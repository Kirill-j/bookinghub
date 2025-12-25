package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"
)

type BookingHandler struct {
	repo    *repo.BookingRepo
	service *service.BookingService
}

func NewBookingHandler(repo *repo.BookingRepo, service *service.BookingService) *BookingHandler {
	return &BookingHandler{repo: repo, service: service}
}

func (h *BookingHandler) My(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}
	items, err := h.repo.ListByUser(r.Context(), uid)
	if err != nil {
		http.Error(w, "Не удалось получить бронирования: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type createBookingReq struct {
	ResourceID uint64 `json:"resourceId"`
	StartAt    string `json:"startAt"` // ISO-строка
	EndAt      string `json:"endAt"`
}

// ожидаем формат RFC3339, например: 2025-12-25T10:00:00
func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Поддержим без timezone (локально) + с timezone
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05", s)
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	uid := GetUserID(r)
	if uid == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	var req createBookingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	startAt, err := parseTime(req.StartAt)
	if err != nil {
		http.Error(w, "Некорректное startAt. Формат: YYYY-MM-DDTHH:MM:SS", http.StatusBadRequest)
		return
	}
	endAt, err := parseTime(req.EndAt)
	if err != nil {
		http.Error(w, "Некорректное endAt. Формат: YYYY-MM-DDTHH:MM:SS", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), uid, req.ResourceID, startAt, endAt)
	if err != nil {
		if err == service.ErrConflict {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

// Менеджерская часть: список ожидающих
func (h *BookingHandler) Pending(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.ListPending(r.Context())
	if err != nil {
		http.Error(w, "Не удалось получить список: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type updateStatusReq struct {
	Status         domain.BookingStatus `json:"status"` // APPROVED / REJECTED
	ManagerComment *string              `json:"managerComment"`
}

func (h *BookingHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	// id берём из query ?id=123 чтобы не городить роут-параметры.
	// Потом улучшим на /bookings/{id}.
	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if idStr == "" {
		http.Error(w, "Нужен параметр id", http.StatusBadRequest)
		return
	}
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	var req updateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if req.Status != domain.BookingApproved && req.Status != domain.BookingRejected {
		http.Error(w, "status должен быть APPROVED или REJECTED", http.StatusBadRequest)
		return
	}

	b, err := h.repo.GetByID(r.Context(), uint64(id64))
	if err != nil {
		http.Error(w, "Ошибка базы: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if b == nil {
		http.Error(w, "Бронирование не найдено", http.StatusNotFound)
		return
	}
	if b.Status != domain.BookingPending {
		http.Error(w, "Можно менять статус только у брони со статусом PENDING", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), uint64(id64), req.Status, req.ManagerComment); err != nil {
		http.Error(w, "Не удалось обновить статус: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
