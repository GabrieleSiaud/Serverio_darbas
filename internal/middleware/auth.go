package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"serverio_darbas/internal/auth"
	"serverio_darbas/internal/generated/repository"

	"go.uber.org/zap"
)

type AuthMiddleware struct {
	authService *auth.AuthService
	log         *zap.SugaredLogger
}

func NewAuthMiddleware(authService *auth.AuthService, log *zap.SugaredLogger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		log:         log,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := m.authenticateRequest(r)
		if err != nil {
			m.log.Debugw("Authentication failed", "error", err, "path", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
	func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			// Jei roles nereikia, tiesiog pritaikom RequireAuth
			return m.RequireAuth(next)
		}
	}
*/
func (m *AuthMiddleware) authenticateRequest(r *http.Request) (*repository.User, error) {
	// Try JWT token first (for API/mobile)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return m.authService.ValidateJWT(r.Context(), token)
	}

	// Try session cookie (for web)
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		return m.authService.ValidateSession(r.Context(), cookie.Value)
	}

	return nil, errors.New("no valid authentication found")
}
