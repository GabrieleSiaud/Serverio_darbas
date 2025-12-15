package router

import (
	"net/http"

	"serverio_darbas/internal/handlers"
	authmw "serverio_darbas/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	gameHandler *handlers.GameHandler,
	reviewHandler *handlers.ReviewHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *authmw.AuthMiddleware,
	externalHandler *handlers.ExternalHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Routes
	r.Route("/user", func(r chi.Router) {
		// ✅ tik admin gali matyti visus users
		r.With(authMiddleware.RequireRole("admin")).Get("/", userHandler.GetUsers)

		// paliekam kaip buvo
		r.Post("/", userHandler.CreateUser)
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", gameHandler.ListGames)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// ✅ šituos verta apsaugoti
		r.With(authMiddleware.RequireAuth).Get("/me", authHandler.Me)
		r.With(authMiddleware.RequireAuth).Post("/logout", authHandler.Logout)
	})
	r.Route("/reviews", func(r chi.Router) {
		r.With(authMiddleware.RequireRole("moderator")).Delete("/{reviewID}", func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "reviewID")
			id, err := uuid.Parse(idStr)
			if err != nil {
				http.Error(w, "invalid reviewID", http.StatusBadRequest)
				return
			}
			reviewHandler.DeleteReview(w, r, id)
		})
	})
	r.Route("/external", func(r chi.Router) {
		r.Get("/deals", externalHandler.Deals)
	})

	r.Route("/auth/battlenet", func(r chi.Router) {
		r.Get("/login", authHandler.BattleNetLogin)
		r.Get("/callback", authHandler.BattleNetCallback)
	})

	return r
}
