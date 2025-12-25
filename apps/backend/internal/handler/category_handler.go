package handler

import (
	"net/http"

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
