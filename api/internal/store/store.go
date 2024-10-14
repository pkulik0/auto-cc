package store

import (
	"context"
	"time"

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
	CreateSessionGoogle(ctx context.Context, userID, accessToken, refreshToken, scopes string, expiry time.Time, credentials model.CredentialsGoogle) (*model.SessionGoogle, error)
	// GetSessionGoogleByCredentialsID returns a Google API session by credentials ID.
	GetSessionGoogleByCredentialsID(ctx context.Context, credentialsID uint, userID string) (*model.SessionGoogle, error)
	// GetSessionGoogleAll returns all Google API sessions for a user.
	GetSessionGoogleAll(ctx context.Context, userID string) ([]model.SessionGoogle, error)
	// GetUserSessionGoogleByQuotaAvailable returns a Google API session with N cost to spend.
	// It updates the credentials usage and returns a function to revert the operation.
	GetSessionGoogleByAvailableCost(ctx context.Context, userID string, cost uint) (*model.SessionGoogle, func() error, error)
	// UpdateSessionGoogle updates a Google API session.
	UpdateSessionGoogle(ctx context.Context, session *model.SessionGoogle) error
	// RemoveSessionGoogle removes a Google API session.
	RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error

	// SaveSessionState saves a state value used in OAuth2.
	SaveSessionState(ctx context.Context, credentialsID uint, userID, state, scopes, redirectURL string) error
	// GetSessionState returns a state value used in OAuth2.
	GetSessionState(ctx context.Context, state string) (*model.SessionState, error)
}
