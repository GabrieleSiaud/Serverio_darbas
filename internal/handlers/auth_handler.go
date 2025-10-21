package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"serverio_darbas/internal/auth"

	"github.com/markbates/goth"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// --- Register ---
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.RegisterUser(r.Context(), req.Email, req.Password, req.Name, req.Surname, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// --- Login ---
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Device info + IP
	deviceInfo := r.UserAgent()
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	authResponse, err := h.authService.LoginWithPassword(r.Context(), req.Email, req.Password, deviceInfo, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    authResponse.SessionToken,
		Path:     "/",
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login successful",
		"user":    authResponse.User,
	})
}

// --- Validate session ---
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "no session", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) BattleNetLogin(w http.ResponseWriter, r *http.Request) {
	provider, _ := goth.GetProvider("battlenet")
	state := "xyz123" // generuok random string
	sess, _ := provider.BeginAuth(state)
	url, _ := sess.GetAuthURL()

	// Marshal + Base64 encode
	encoded := base64.StdEncoding.EncodeToString([]byte(sess.Marshal()))

	http.SetCookie(w, &http.Cookie{
		Name:  "battle_session",
		Value: encoded,
		Path:  "/",
	})

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) BattleNetCallback(w http.ResponseWriter, r *http.Request) {
	provider, _ := goth.GetProvider("battlenet")

	cookie, err := r.Cookie("battle_session")
	if err != nil {
		http.Error(w, "missing session cookie", http.StatusBadRequest)
		return
	}

	// Base64 decode
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		http.Error(w, "decode cookie error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Atstatom sesiją
	sess, err := provider.UnmarshalSession(string(data))
	if err != nil {
		http.Error(w, "unmarshal session error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := sess.Authorize(provider, r.URL.Query()); err != nil {
		http.Error(w, "authorize error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := provider.FetchUser(sess)
	if err != nil {
		http.Error(w, "fetch user error: "+err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// --- Logout ---
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "no session", http.StatusUnauthorized)
		return
	}

	if err := h.authService.Logout(r.Context(), cookie.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "logged out",
	})
}
