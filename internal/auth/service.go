package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	//"net"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth"
	"serverio_darbas/internal/generated/repository"
)

type AuthService struct {
	queries    *repository.Queries
	jwtService *JWTService
}

type AuthResponse struct {
	User         *repository.User `json:"user"`
	AccessToken  string           `json:"access_token,omitempty"`
	SessionToken string           `json:"session_token,omitempty"`
}

// API-specific auth response for mobile/API clients
type APIAuthResponse struct {
	Success     bool             `json:"success"`
	User        *repository.User `json:"user"`
	AccessToken string           `json:"access_token"`
	TokenType   string           `json:"token_type"`
	ExpiresIn   int              `json:"expires_in"`
	Provider    string           `json:"provider,omitempty"`
}

func NewAuthService(queries *repository.Queries, jwtService *JWTService) *AuthService {
	return &AuthService{
		queries:    queries,
		jwtService: jwtService,
	}
}

// Get all OAuth providers for a user
func (a *AuthService) GetUserOAuthProviders(ctx context.Context, userID uuid.UUID) ([]repository.OauthProvider, error) {
	return a.queries.GetUserOAuthProviders(ctx, userID)
}

// Handle API OAuth callback
func (a *AuthService) HandleAPIAuthCallback(ctx context.Context, gothUser goth.User, deviceInfo, ipAddress string) (*APIAuthResponse, error) {
	authResponse, err := a.HandleOAuthCallback(ctx, gothUser, deviceInfo, ipAddress)
	if err != nil {
		return nil, err
	}

	return &APIAuthResponse{
		Success:     true,
		User:        authResponse.User,
		AccessToken: authResponse.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours
		Provider:    gothUser.Provider,
	}, nil
}

// Login with email + password
func (a *AuthService) LoginWithPassword(ctx context.Context, email, password, deviceInfo, ipAddress string) (*AuthResponse, error) {
	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !CheckPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return a.createUserSession(ctx, &user, deviceInfo, ipAddress)
}

// Register new user
func (a *AuthService) RegisterUser(ctx context.Context, email, password, name, surname, username string) (*repository.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	params := repository.CreateUserParams{
		Name:     name,
		Surname:  surname,
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	user, err := a.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Handle OAuth callback
func (a *AuthService) HandleOAuthCallback(ctx context.Context, gothUser goth.User, deviceInfo, ipAddress string) (*AuthResponse, error) {
	// Check if OAuth provider exists
	oauthProvider, err := a.queries.GetOAuthProviderByExternalID(ctx, repository.GetOAuthProviderByExternalIDParams{
		Provider:       gothUser.Provider,
		ProviderUserID: gothUser.UserID,
	})

	var user *repository.User

	if err == nil {
		// Existing OAuth user
		existingUser, err := a.queries.GetUserByID(ctx, oauthProvider.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		user = &existingUser

		// Update tokens

		/*var accessToken, refreshToken *string
		if gothUser.AccessToken != "" {
			accessToken = &gothUser.AccessToken
		}
		if gothUser.RefreshToken != "" {
			refreshToken = &gothUser.RefreshToken
		}*/

		_, err = a.queries.LinkOAuthProvider(ctx, repository.LinkOAuthProviderParams{
			UserID:           user.ID,
			Provider:         gothUser.Provider,
			ProviderUserID:   gothUser.UserID,
			ProviderUsername: pgtype.Text{String: gothUser.NickName, Valid: gothUser.NickName != ""},
			ProviderEmail:    pgtype.Text{String: gothUser.Email, Valid: gothUser.Email != ""},
			AccessToken:      pgtype.Text{String: gothUser.AccessToken, Valid: gothUser.AccessToken != ""},
			RefreshToken:     pgtype.Text{String: gothUser.RefreshToken, Valid: gothUser.RefreshToken != ""},
			TokenExpiresAt:   TimeToPgTimestamptz(gothUser.ExpiresAt),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to update OAuth tokens: %w", err)
		}
	} else {
		// New OAuth user
		email := gothUser.Email
		if email == "" {
			email = fmt.Sprintf("%s@%s.oauth", gothUser.UserID, gothUser.Provider)
		}

		name := gothUser.FirstName
		if name == "" {
			name = gothUser.NickName
		}
		if name == "" {
			name = "User"
		}

		surname := gothUser.LastName
		username := gothUser.NickName

		newUser, err := a.RegisterUser(ctx, email, "", name, surname, username)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		user = newUser

		// Link OAuth provider
		_, err = a.queries.LinkOAuthProvider(ctx, repository.LinkOAuthProviderParams{
			UserID:           user.ID,
			Provider:         gothUser.Provider,
			ProviderUserID:   gothUser.UserID,
			ProviderUsername: pgtype.Text{String: gothUser.NickName, Valid: true},
			ProviderEmail:    pgtype.Text{String: gothUser.Email, Valid: true},
			AccessToken:      pgtype.Text{String: gothUser.AccessToken, Valid: true},
			RefreshToken:     pgtype.Text{String: gothUser.RefreshToken, Valid: true},
			TokenExpiresAt:   TimeToPgTimestamptz(gothUser.ExpiresAt),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to link OAuth provider: %w", err)
		}
	}

	// Update last login
	err = a.queries.UpdateUserLastLogin(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	return a.createUserSession(ctx, user, deviceInfo, ipAddress)
}

// Create a session
func (a *AuthService) createUserSession(ctx context.Context, user *repository.User, deviceInfo, ipAddress string) (*AuthResponse, error) {
	sessionToken, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	// JWT generation
	jwtToken, jti, err := a.jwtService.GenerateToken(user.ID, user.Email, "user")
	if err != nil {
		return nil, err
	}

	// konvertuojam IP į netip.Addr pointerį
	var ipAddr *netip.Addr
	if ipAddress != "" {
		a, err := netip.ParseAddr(ipAddress)
		if err == nil {
			ipAddr = &a
		}
	}

	// DeviceInfo į pgtype.Text
	deviceInfoText := pgtype.Text{
		String: deviceInfo,
		Valid:  deviceInfo != "",
	}

	// CreateSession
	_, err = a.queries.CreateSession(ctx, repository.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		JwtTokenID:   pgtype.Text{String: jti, Valid: true},
		DeviceInfo:   deviceInfoText,
		IpAddress:    ipAddr,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  jwtToken,
		SessionToken: sessionToken,
	}, nil
}

// Validate session token
func (a *AuthService) ValidateSession(ctx context.Context, sessionToken string) (*repository.User, error) {
	// gauname session
	sessionData, err := a.queries.GetSessionByToken(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	// atnaujinam last_used timestamp
	_, _ = a.queries.UpdateSessionLastUsed(ctx, sessionData.SessionID)

	// gauname User info pagal sessionData.UserID
	user, err := a.queries.GetUserByID(ctx, sessionData.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// slaptažodžio negrąžinam
	user.Password = ""

	return &user, nil
}

// Validate JWT
func (a *AuthService) ValidateJWT(ctx context.Context, tokenString string) (*repository.User, error) {
	claims, err := a.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Name:  "",
	}

	return user, nil
}

// Logout
func (a *AuthService) Logout(ctx context.Context, sessionToken string) error {
	return a.queries.DeleteSessionByToken(ctx, sessionToken)
}

// Helper to generate random token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
