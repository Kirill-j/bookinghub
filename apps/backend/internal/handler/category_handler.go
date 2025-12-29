package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"bookinghub-backend/internal/repo"
)

type CategoryHandler struct {
	repo *repo.CategoryRepo
}

func NewCategoryHandler(repo *repo.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list categories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type createCategoryReq struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		http.Error(w, "name обязателен", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(r.Context(), req.Name)
	if err != nil {
		http.Error(w, "failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	var req createCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		http.Error(w, "name обязателен", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), uint64(id64), name); err != nil {
		http.Error(w, "failed to update category: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), uint64(id64)); err != nil {
		// Частый кейс: на категорию есть ресурсы → FK не даст удалить
		http.Error(w, "Не удалось удалить категорию (возможно, к ней привязаны объявления): "+err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
