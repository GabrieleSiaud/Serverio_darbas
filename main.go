package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
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

	r.Get("/reviews/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		gameIDStr := chi.URLParam(r, "id")
		gameID, err := uuid.Parse(gameIDStr)
		if err != nil {
			http.Error(w, "invalid game id", http.StatusBadRequest)
			return
		}

		reviews, err := queries.ListReviewsByGame(ctx, gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, rev := range reviews {
			comment := ""
			if rev.Comment.Valid {
				comment = rev.Comment.String
			}
			created := ""
			if rev.CreatedAt.Valid {
				created = rev.CreatedAt.Time.String()
			}
			fmt.Fprintf(w, "User: %s | Rating: %d | Comment: %s | Created: %s\n",
				rev.Username, rev.Rating, comment, created)
		}
	})

	// GET all saved games for a user
	r.Get("/saved-games/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIDStr := chi.URLParam(r, "id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		savedGames, err := queries.ListSavedGamesByUser(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, sg := range savedGames {
			desc := ""
			if sg.Description.Valid {
				desc = sg.Description.String
			}
			release := ""
			if sg.ReleaseDate.Valid {
				release = sg.ReleaseDate.Time.String()
			}
			created := ""
			if sg.CreatedAt.Valid {
				created = sg.CreatedAt.Time.String()
			}
			fmt.Fprintf(w, "Game: %s | Description: %s | Release: %s | SavedAt: %s\n",
				sg.Title, desc, release, created)
		}
	})

	// POST to save a game for a user
	r.Post("/saved-games", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req struct {
			UserID uuid.UUID `json:"user_id"`
			GameID uuid.UUID `json:"game_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		savedGame, err := queries.SaveGame(ctx, repository.SaveGameParams{
			UserID: req.UserID,
			GameID: req.GameID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Saved Game ID: %s for User: %s and Game: %s\n",
			savedGame.ID, savedGame.UserID, savedGame.GameID)
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
