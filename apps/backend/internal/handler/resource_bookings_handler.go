package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"bookinghub-backend/internal/repo"
)

type ResourceBookingsHandler struct {
	bookings *repo.BookingRepo
}

func NewResourceBookingsHandler(bookings *repo.BookingRepo) *ResourceBookingsHandler {
	return &ResourceBookingsHandler{bookings: bookings}
}

func (h *ResourceBookingsHandler) List(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		http.Error(w, "Некорректный id ресурса", http.StatusBadRequest)
		return
	}

	fromStr := strings.TrimSpace(r.URL.Query().Get("from"))
	toStr := strings.TrimSpace(r.URL.Query().Get("to"))
	if fromStr == "" || toStr == "" {
		http.Error(w, "Нужны параметры from и to в формате YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		http.Error(w, "Некорректный from", http.StatusBadRequest)
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		http.Error(w, "Некорректный to", http.StatusBadRequest)
		return
	}

	// to делаем эксклюзивным (следующий день 00:00)
	to = to.Add(24 * time.Hour)

	items, err := h.bookings.ListByResourceBetween(r.Context(), uint64(id64), from, to)
	if err != nil {
		http.Error(w, "Не удалось получить бронирования: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
