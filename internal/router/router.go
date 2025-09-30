package router

import (
	"net/http"
	"serverio_darbas/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(userHandler *handlers.UserHandler, gameHandler *handlers.GameHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // papildomas patarimas: sugauti panics

	// Routes
	r.Route("/user", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", gameHandler.ListGames)
	})

	return r
}
