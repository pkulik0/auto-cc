package store

import (
	"context"

	"github.com/pkulik0/autocc/api/internal/model"
)

// Store is an interface for storing external services and client credentials.
//
//go:generate mockgen -destination=../mock/store.go -package=mock . Store
type Store interface {
	// Transaction executes store operations in a transaction.
	Transaction(ctx context.Context, f func(ctx context.Context, store Store) error) error

	// AddCredentialsGoogle adds Google client credentials to the store.
	AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error)
	// AddCredentialsDeepL adds DeepL client credentials to the store.
	AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error)

	// GetCredentialsGoogleAll returns all Google client credentials.
	GetCredentialsGoogleAll(ctx context.Context) ([]model.CredentialsGoogle, error)
	// GetCredentialsDeepLAll returns all DeepL client credentials.
	GetCredentialsDeepLAll(ctx context.Context) ([]model.CredentialsDeepL, error)

	// GetCredentialsGoogleByID returns Google client credentials by ID.
	GetCredentialsGoogleByID(ctx context.Context, id uint) (*model.CredentialsGoogle, error)
	// GetCredentialsDeepLByID returns DeepL client credentials by ID.
	GetCredentialsDeepLByID(ctx context.Context, id uint) (*model.CredentialsDeepL, error)

	// RemoveCredentialsGoogle removes Google client credentials from the store.
	RemoveCredentialsGoogle(ctx context.Context, id uint) error
	// RemoveCredentialsDeepL removes DeepL client credentials from the store.
	RemoveCredentialsDeepL(ctx context.Context, id uint) error

	// CreateSessionGoogle creates a new Google API session.
	CreateSessionGoogle(ctx context.Context, userID, accessToken, refreshToken string, expiry int64, credentialsID uint, scopes string) (*model.SessionGoogle, error)
	// GetUserSessionsGoogle returns all Google API sessions for a user.
	GetUserSessionsGoogle(ctx context.Context, userID string) ([]model.SessionGoogle, error)
	// RemoveSessionGoogle removes a Google API session.
	RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error

	// SaveSessionState saves a state value used in OAuth2.
	SaveSessionState(ctx context.Context, credentialsID uint, userID, state string, scopes string) error
	// GetSessionState returns a state value used in OAuth2.
	GetSessionState(ctx context.Context, state string) (*model.SessionState, error)
}
