package services

import (
	"context"
	"fmt"
)

func TestUserAccessToken(d *Diablo3Service, ctx context.Context, userID string) {
	err := d.TestAccessToken(ctx, userID)
	if err != nil {
		fmt.Printf("❌ Access token is invalid or expired: %v\n", err)
	} else {
		fmt.Println("✅ Access token is valid!")
	}
}

// Galima testuoti main funkcijoje
func RunAccessTokenTest(d *Diablo3Service, userID string) {
	ctx := context.Background()
	TestUserAccessToken(d, ctx, userID)
}
