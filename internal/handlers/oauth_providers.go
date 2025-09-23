package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"serverio_darbas/internal/generated/repository"
)

type OAuthProviderHandler struct {
	db *repository.Queries
}

func NewOAuthProviderHandler(db *repository.Queries) *OAuthProviderHandler {
	return &OAuthProviderHandler{db: db}
}

// LinkOAuthProviderRequest struktūra POST /link
type LinkOAuthProviderRequest struct {
	UserID           uuid.UUID `json:"user_id"`
	Provider         string    `json:"provider"`
	ProviderUserID   string    `json:"provider_user_id"`
	ProviderUsername string    `json:"provider_username"`
	ProviderEmail    string    `json:"provider_email"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenExpiresAt   time.Time `json:"token_expires_at"`
}

// LinkOAuthProvider sukuria arba atnaujina OAuth providerį
func (h *OAuthProviderHandler) LinkOAuthProvider(w http.ResponseWriter, r *http.Request) {
	var req LinkOAuthProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	provider, err := h.db.LinkOAuthProvider(r.Context(), repository.LinkOAuthProviderParams{
		UserID:           req.UserID,
		Provider:         req.Provider,
		ProviderUserID:   req.ProviderUserID,
		ProviderUsername: pgtype.Text{String: req.ProviderUsername, Valid: true},
		ProviderEmail:    pgtype.Text{String: req.ProviderEmail, Valid: true},
		AccessToken:      pgtype.Text{String: req.AccessToken, Valid: true},
		RefreshToken:     pgtype.Text{String: req.RefreshToken, Valid: true},
		TokenExpiresAt:   pgtype.Timestamptz{Time: req.TokenExpiresAt, Valid: true},
	})
	if err != nil {
		http.Error(w, "failed to link oauth provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

// GetOAuthProviderByExternalIDRequest struktūra GET /provider/external
type GetOAuthProviderByExternalIDRequest struct {
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
}

func (h *OAuthProviderHandler) GetOAuthProviderByExternalID(w http.ResponseWriter, r *http.Request, provider, providerUserID string) {
	oauth, err := h.db.GetOAuthProviderByExternalID(r.Context(), repository.GetOAuthProviderByExternalIDParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if err != nil {
		http.Error(w, "oauth provider not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oauth)
}

// GetOAuthProviderByUserRequest struktūra GET /provider/user
func (h *OAuthProviderHandler) GetOAuthProviderByUser(w http.ResponseWriter, r *http.Request, userID uuid.UUID, provider string) {
	oauth, err := h.db.GetOAuthProviderByUser(r.Context(), repository.GetOAuthProviderByUserParams{
		UserID:   userID,
		Provider: provider,
	})
	if err != nil {
		http.Error(w, "oauth provider not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oauth)
}

// DeleteOAuthProviderRequest struktūra DELETE /provider
func (h *OAuthProviderHandler) DeleteOAuthProvider(w http.ResponseWriter, r *http.Request, userID uuid.UUID, provider string) {
	if err := h.db.DeleteOAuthProvider(r.Context(), repository.DeleteOAuthProviderParams{
		UserID:   userID,
		Provider: provider,
	}); err != nil {
		http.Error(w, "failed to delete oauth provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "oauth provider deleted",
	})
}
