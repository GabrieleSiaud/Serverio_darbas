package handlers

import (
	"encoding/json"
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"serverio_darbas/internal/generated/repository"
)

type UserSessionHandler struct {
	db *repository.Queries
}

func NewUserSessionHandler(db *repository.Queries) *UserSessionHandler {
	return &UserSessionHandler{db: db}
}

type CreateSessionRequest struct {
	UserID       uuid.UUID   `json:"user_id"`
	SessionToken string      `json:"session_token"`
	JWTTokenID   string      `json:"jwt_token_id"`
	DeviceInfo   string      `json:"device_info"`
	IPAddress    *netip.Addr `json:"ip_address"`
	ExpiresAt    time.Time   `json:"expires_at"`
}

// CreateSession sukuria naują vartotojo sesiją
func (h *UserSessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.db.CreateSession(r.Context(), repository.CreateSessionParams{
		UserID:       req.UserID,
		SessionToken: req.SessionToken,
		JwtTokenID:   pgtype.Text{String: req.JWTTokenID, Valid: true},
		DeviceInfo:   pgtype.Text{String: req.DeviceInfo, Valid: true},
		IpAddress:    req.IPAddress,
		ExpiresAt:    pgtype.Timestamptz{Time: req.ExpiresAt, Valid: true},
	})

	if err != nil {
		http.Error(w, "failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetSessionByToken grąžina sesiją pagal token
func (h *UserSessionHandler) GetSessionByToken(w http.ResponseWriter, r *http.Request, token string) {
	session, err := h.db.GetSessionByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "failed to fetch session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// DeleteSession ištrina sesiją pagal ID
func (h *UserSessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request, sessionID uuid.UUID) {
	if err := h.db.DeleteSession(r.Context(), sessionID); err != nil {
		http.Error(w, "failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "session deleted",
	})
}

// DeleteExpiredSessions ištrina pasibaigusias sesijas
func (h *UserSessionHandler) DeleteExpiredSessions(w http.ResponseWriter, r *http.Request) {
	if err := h.db.DeleteExpiredSessions(r.Context()); err != nil {
		http.Error(w, "failed to delete expired sessions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "expired sessions deleted",
	})
}
