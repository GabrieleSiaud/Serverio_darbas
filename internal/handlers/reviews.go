package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"serverio_darbas/internal/generated/repository"
)

type ReviewHandler struct {
	db *repository.Queries
}

func NewReviewHandler(db *repository.Queries) *ReviewHandler {
	return &ReviewHandler{db: db}
}

type CreateReviewRequest struct {
	GameID  uuid.UUID `json:"game_id"`
	UserID  uuid.UUID `json:"user_id"`
	Rating  int       `json:"rating"`
	Comment string    `json:"comment"`
}

// CreateReview creates or updates a review
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.db.CreateReview(r.Context(), repository.CreateReviewParams{
		GameID: req.GameID,
		UserID: req.UserID,
		Rating: int16(req.Rating), // tik int32, ne pgtype.Int4
		Comment: pgtype.Text{ // pgtype.Text, ne string
			String: req.Comment,
			Valid:  true,
		},
	})
	if err != nil {
		http.Error(w, "failed to create review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// GetReview fetches a single review by ID
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request, reviewID uuid.UUID) {
	review, err := h.db.GetReview(r.Context(), reviewID)
	if err != nil {
		http.Error(w, "failed to fetch review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// ListReviewsByGame returns all reviews for a game
func (h *ReviewHandler) ListReviewsByGame(w http.ResponseWriter, r *http.Request, gameID uuid.UUID) {
	reviews, err := h.db.ListReviewsByGame(r.Context(), gameID)
	if err != nil {
		http.Error(w, "failed to fetch reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// ListReviewsByUser returns all reviews by a user
func (h *ReviewHandler) ListReviewsByUser(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	reviews, err := h.db.ListReviewsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// DeleteReview removes a review by ID
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request, reviewID uuid.UUID) {
	if err := h.db.DeleteReview(r.Context(), reviewID); err != nil {
		http.Error(w, "failed to delete review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "review deleted",
	})
}
