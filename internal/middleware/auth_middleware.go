package middleware

import (
	"net/http"
	"serverio_darbas/internal/generated/repository"
)

func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userAny := r.Context().Value("user")
			user, ok := userAny.(*repository.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			roles, err := m.authService.GetUserRoles(r.Context(), user.ID)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, rn := range roles {
				if rn == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		}))
	}
}
