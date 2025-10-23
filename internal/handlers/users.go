package handlers

import (
	"encoding/json"
	"net/http"
	"serverio_darbas/internal/auth"

	"serverio_darbas/internal/generated/repository"
)

type UserHandler struct {
	db *repository.Queries
}

func NewUserHandler(db *repository.Queries) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.db.ListUsers(ctx)
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req repository.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 👇 Hashinam slaptažodį prieš įrašymą į DB
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	req.Password = hashedPassword

	// Sukuriam vartotoją su hashuotu passwordu
	user, err := h.db.CreateUser(ctx, req)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
