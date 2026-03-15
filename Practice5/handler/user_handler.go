package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"practice5/models"
	"practice5/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GET /users
// Query params: page, page_size, order_by, order_dir, id, name, email, gender, birthdate
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	params := models.FilterParams{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  q.Get("order_by"),
		OrderDir: q.Get("order_dir"),
	}

	if v := q.Get("id"); v != "" {
		id, err := strconv.Atoi(v)
		if err == nil {
			params.ID = &id
		}
	}
	if v := q.Get("name"); v != "" {
		params.Name = &v
	}
	if v := q.Get("email"); v != "" {
		params.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		params.Gender = &v
	}
	if v := q.Get("birthdate"); v != "" {
		params.Birthdate = &v
	}

	result, err := h.repo.GetPaginatedUsers(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}


func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	user1, err1 := strconv.Atoi(q.Get("user1"))
	user2, err2 := strconv.Atoi(q.Get("user2"))

	if err1 != nil || err2 != nil {
		http.Error(w, "user1 and user2 must be valid integers", http.StatusBadRequest)
		return
	}
	if user1 == user2 {
		http.Error(w, "user1 and user2 must be different", http.StatusBadRequest)
		return
	}

	friends, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
