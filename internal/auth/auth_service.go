package auth

import (
	"context"

	"github.com/google/uuid"
)

func (s *AuthService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.queries.ListUserRoles(ctx, userID)
}
