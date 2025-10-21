package router

import (
	"net/http"
	"serverio_darbas/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Pridėtas authHandler kaip trečias parametras
func NewRouter(userHandler *handlers.UserHandler, gameHandler *handlers.GameHandler, authHandler *handlers.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/user", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", gameHandler.ListGames)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/me", authHandler.Me)
		r.Post("/logout", authHandler.Logout)
	})

	r.Route("/auth/battlenet", func(r chi.Router) {
		r.Get("/login", authHandler.BattleNetLogin)
		r.Get("/callback", authHandler.BattleNetCallback)
	})

	return r
}
