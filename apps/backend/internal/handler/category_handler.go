package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
