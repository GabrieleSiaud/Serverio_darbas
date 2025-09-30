package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"serverio_darbas/internal/generated/repository"
	"serverio_darbas/internal/handlers"
	"serverio_darbas/internal/router"

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

	queries := repository.New(pool)
	userHandler := handlers.NewUserHandler(queries)
	gameHandler := handlers.NewGameHandler(queries)

	r := router.NewRouter(userHandler, gameHandler)

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
