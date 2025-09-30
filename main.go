package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"serverio_darbas/internal/generated/repository"
	"serverio_darbas/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// DB URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:root@localhost:5432/serverio_duomenubaze?sslmode=disable"
	}

	// Connect to DB
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}
	fmt.Println("✅ Database connected!")

	// Init SQLC repository
	queries := repository.New(pool)

	// Init UserHandler
	userHandler := handlers.NewUserHandler(queries)
	gameHandler := handlers.NewGameHandler(queries)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/user", userHandler.GetUsers)
	r.Post("/user", userHandler.CreateUser)
	r.Get("/games", gameHandler.ListGames)

	// Start server
	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fmt.Println("🚀 Server running on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
