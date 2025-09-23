package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"serverio_darbas/internal/generated/repository"
)

type SavedGameHandler struct {
	db *repository.Queries
}

func NewSavedGameHandler(db *repository.Queries) *SavedGameHandler {
	return &SavedGameHandler{db: db}
}

// POST /saved_games
func (h *SavedGameHandler) SaveGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID string `json:"user_id"`
		GameID string `json:"game_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	gameID, err := uuid.Parse(req.GameID)
	if err != nil {
		http.Error(w, "invalid game_id", http.StatusBadRequest)
		return
	}

	saved, err := h.db.SaveGame(ctx, repository.SaveGameParams{
		UserID: userID,
		GameID: gameID,
	})
	if err != nil {
		http.Error(w, "failed to save game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(saved)
}

// GET /users/{userID}/saved_games
func (h *SavedGameHandler) ListSavedGamesByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	games, err := h.db.ListSavedGamesByUser(ctx, userID)
	if err != nil {
		http.Error(w, "failed to fetch saved games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(games)
}

// DELETE /users/{userID}/saved_games/{gameID}
func (h *SavedGameHandler) DeleteSavedGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := chi.URLParam(r, "userID")
	gameIDStr := chi.URLParam(r, "gameID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "invalid game_id", http.StatusBadRequest)
		return
	}

	err = h.db.DeleteSavedGame(ctx, repository.DeleteSavedGameParams{
		UserID: userID,
		GameID: gameID,
	})
	if err != nil {
		http.Error(w, "failed to delete saved game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
