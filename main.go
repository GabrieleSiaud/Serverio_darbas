package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"serverio_darbas/internal/seed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"serverio_darbas/internal/generated/repository"
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

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if err := seed.Seed(queries); err != nil {
		log.Fatal(err)
	}

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Users endpoint
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		users, err := queries.ListUsers(ctx) // <- turi būti SQLC generuotas metodas
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, u := range users {
			fmt.Fprintf(w, "%s %s - %s\n", u.Name, u.Surname, u.Email)
		}
	})

	// Games endpoint
	r.Get("/games", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		games, err := queries.ListGames(ctx) // <- SQLC generuotas metodas tavo queries
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, g := range games {
			fmt.Fprintf(w, "%s - %s\n", g.Title, g.Description)
		}
	})

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
