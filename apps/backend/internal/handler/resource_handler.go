package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"bookinghub-backend/internal/repo"
)

type ResourceHandler struct {
	repo *repo.ResourceRepo
}

func NewResourceHandler(repo *repo.ResourceRepo) *ResourceHandler {
	return &ResourceHandler{repo: repo}
}

func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list resources: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type createResourceRequest struct {
	CategoryID   uint64  `json:"categoryId"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	Location     *string `json:"location"`
	PricePerHour int     `json:"pricePerHour"`
}

func (h *ResourceHandler) My(w http.ResponseWriter, r *http.Request) {
	ownerID := GetUserID(r)
	if ownerID == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	items, err := h.repo.ListByOwner(r.Context(), ownerID)
	if err != nil {
		http.Error(w, "failed to list my resources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.CategoryID == 0 {
		http.Error(w, "categoryId is required", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if req.PricePerHour < 0 {
		http.Error(w, "pricePerHour must be >= 0", http.StatusBadRequest)
		return
	}

	ownerID := GetUserID(r)
	if ownerID == 0 {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	id, err := h.repo.Create(
		r.Context(),
		ownerID,
		req.CategoryID,
		req.Title,
		req.Description,
		req.Location,
		req.PricePerHour,
	)
	if err != nil {
		http.Error(w, "failed to create resource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
