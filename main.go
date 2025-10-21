package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"serverio_darbas/internal/auth"
	"serverio_darbas/internal/generated/repository"
	"serverio_darbas/internal/handlers"
	"serverio_darbas/internal/router"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	fmt.Println("Client ID:", os.Getenv("BATTLE_CLIENT_ID"))
	fmt.Println("Client Secret:", os.Getenv("BATTLE_CLIENT_SECRET"))
	// DB URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:root@localhost:5432/serverio_duomenubaze?sslmode=disable"
	}

	// DB connection
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	defer pool.Close()
	fmt.Println("✅ Database connected!")

	queries := repository.New(pool)

	// Logger
	_, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// JWT service
	jwtService := auth.NewJWTService("slaptas_raktas", string(24*time.Hour)) // Pakeisk į savo slaptą raktą
	authService := auth.NewAuthService(queries, jwtService)

	// Handlers
	userHandler := handlers.NewUserHandler(queries)
	gameHandler := handlers.NewGameHandler(queries)
	authHandler := handlers.NewAuthHandler(authService)
	auth.InitBattleNetOAuth()

	// Router
	r := router.NewRouter(userHandler, gameHandler, authHandler)

	// Start server
	port := "3000"
	fmt.Println("🚀 Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
