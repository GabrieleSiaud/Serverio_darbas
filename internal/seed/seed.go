package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"serverio_darbas/internal/generated/repository"
)

func Seed(db *repository.Queries) error {
	ctx := context.Background()

	// Sukuri vartotoją
	if _, err := db.CreateUser(ctx, repository.CreateUserParams{
		Name:     "Admin",
		Surname:  "Admin",
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
	}); err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}

	// Sukuri žaidimą
	if _, err := db.CreateGame(ctx, repository.CreateGameParams{
		Title: "Diablo III",
		Description: pgtype.Text{
			String: "Action RPG",
			Valid:  true,
		},
		ReleaseDate: pgtype.Date{
			Time:  time.Date(2012, 5, 15, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("cannot create game: %w", err)
	}

	return nil
}
