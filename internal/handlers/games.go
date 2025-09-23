package handlers

import (
	"encoding/json"
	"net/http"

	"serverio_darbas/internal/generated/repository"
)

type GameHandler struct {
	db *repository.Queries
}

func NewGameHandler(db *repository.Queries) *GameHandler {
	return &GameHandler{db: db}
}

// GET /games
func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	games, err := h.db.ListGames(ctx)
	if err != nil {
		http.Error(w, "failed to fetch games", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(games)
}
