package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"serverio_darbas/internal/services"
	"time"

	"serverio_darbas/internal/auth"
	"serverio_darbas/internal/generated/repository"
	"serverio_darbas/internal/handlers"
	middleware "serverio_darbas/internal/middleware"
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
	// Logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	// JWT service
	jwtService := auth.NewJWTService("slaptas_raktas", string(24*time.Hour)) // Pakeisk į savo slaptą raktą
	authService := auth.NewAuthService(queries, jwtService)
	authMw := middleware.NewAuthMiddleware(authService, sugar)

	ctx := context.Background()
	admin, err := queries.GetUserByEmail(ctx, "admin@example.com")
	if err == nil {
		if len(admin.Password) < 60 { // bcrypt hash visada ilgas
			hashed, err := auth.HashPassword(admin.Password)
			if err != nil {
				log.Println("❌ Nepavyko suhashinti admin slaptažodžio:", err)
			} else {
				_, err := pool.Exec(ctx, `UPDATE users SET password=$1 WHERE email=$2`, hashed, admin.Email)
				if err != nil {
					log.Println("❌ Nepavyko atnaujinti admin slaptažodžio:", err)
				} else {
					log.Println("✅ Admin slaptažodis automatiškai suhashintas!")
				}
			}
		}
	}

	// Handlers
	userHandler := handlers.NewUserHandler(queries)
	gameHandler := handlers.NewGameHandler(queries)
	authHandler := handlers.NewAuthHandler(authService)
	reviewHandler := handlers.NewReviewHandler(queries)
	dealsService := services.NewDealsService(queries)
	externalHandler := handlers.NewExternalHandler(dealsService)
	auth.InitBattleNetOAuth()

	// Router
	r := router.NewRouter(userHandler, gameHandler, reviewHandler, authHandler, authMw, externalHandler)

	// Start server
	port := "3000"
	fmt.Println("🚀 Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
