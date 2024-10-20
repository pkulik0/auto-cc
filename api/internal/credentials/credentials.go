package credentials

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/model"
	"github.com/pkulik0/autocc/api/internal/oauth"
	"github.com/pkulik0/autocc/api/internal/store"
	"github.com/pkulik0/autocc/api/internal/translation"
)

// Credentials is an interface for the credentials service.
//
//go:generate mockgen -destination=../mock/credentials.go -package=mock . Credentials
type Credentials interface {
	// AddCredentialsGoogle adds Google client credentials to the store.
	AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error)
	// AddCredentialsDeepL adds DeepL client credentials to the store.
	AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error)

	// GetCredentials returns all client credentials.
	GetCredentials(ctx context.Context) ([]model.CredentialsGoogle, []model.CredentialsDeepL, error)

	// RemoveCredentialsGoogle removes Google client credentials from the store.
	RemoveCredentialsGoogle(ctx context.Context, id uint) error
	// RemoveCredentialsDeepL removes DeepL client credentials from the store.
	RemoveCredentialsDeepL(ctx context.Context, id uint) error

	// GetSessionGoogleURL returns a URL for authenticating with Google.
	GetSessionGoogleURL(ctx context.Context, credentialsID uint, userID, redirectURL string) (string, error)
	// CreateSessionGoogle creates a new Google API session.
	CreateSessionGoogle(ctx context.Context, state, code string) (string, error)
	// RemoveSessionGoogle removes a Google API session.
	RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error
	// GetSessionsGoogleByUser returns all Google API sessions for a user.
	GetSessionsGoogleByUser(ctx context.Context, userID string) ([]model.SessionGoogle, error)
}

var _ Credentials = &credentials{}

type credentials struct {
	store       store.Store
	oauth       oauth.Configs
	translation translation.Translator
}

// New creates a new credentials service.
func New(s store.Store, o oauth.Configs, t translation.Translator) *credentials {
	log.Debug().Msg("created credentials service")
	return &credentials{
		store:       s,
		oauth:       o,
		translation: t,
	}
}

func (c *credentials) AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errs.InvalidInput
	}

	return c.store.AddCredentialsGoogle(ctx, clientID, clientSecret)
}

func (c *credentials) AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error) {
	if key == "" {
		return nil, errs.InvalidInput
	}

	usage, err := c.translation.GetUsageDeepL(ctx, key)
	if err != nil {
		return nil, err
	}

	return c.store.AddCredentialsDeepL(ctx, key, usage)
}

func (c *credentials) GetCredentials(ctx context.Context) ([]model.CredentialsGoogle, []model.CredentialsDeepL, error) {
	google, err := c.store.GetCredentialsGoogleAll(ctx)
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		google = []model.CredentialsGoogle{}
	default:
		return nil, nil, err
	}

	deepL, err := c.store.GetCredentialsDeepLAll(ctx)
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		deepL = []model.CredentialsDeepL{}
	default:
		return nil, nil, err
	}

	return google, deepL, nil
}

func (c *credentials) RemoveCredentialsGoogle(ctx context.Context, id uint) error {
	return c.store.RemoveCredentialsGoogle(ctx, id)
}

func (c *credentials) RemoveCredentialsDeepL(ctx context.Context, id uint) error {
	return c.store.RemoveCredentialsDeepL(ctx, id)
}

func generateState() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

var (
	urlRegex = regexp.MustCompile(`^https?://`)
)

func (c *credentials) GetSessionGoogleURL(ctx context.Context, credentialsID uint, userID, redirectURL string) (string, error) {
	if userID == "" {
		return "", errs.InvalidInput
	}
	if !urlRegex.MatchString(redirectURL) {
		return "", errs.InvalidInput
	}

	credentials, err := c.store.GetCredentialsGoogleByID(ctx, credentialsID)
	if err != nil {
		return "", err
	}
	oauthClient, scopes := c.oauth.GetGoogle(credentials.ClientID, credentials.ClientSecret)

	state, err := generateState()
	if err != nil {
		return "", err
	}
	err = c.store.SaveSessionState(ctx, credentialsID, userID, state, scopes, redirectURL)
	if err != nil {
		return "", err
	}

	return oauthClient.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (c *credentials) CreateSessionGoogle(ctx context.Context, state, code string) (string, error) {
	if state == "" || code == "" {
		return "", errs.InvalidInput
	}

	sessionState, err := c.store.GetSessionState(ctx, state)
	if err != nil {
		return "", err
	}

	credentials, err := c.store.GetCredentialsGoogleByID(ctx, sessionState.CredentialsID)
	if err != nil {
		return "", err
	}
	oauthClient, _ := c.oauth.GetGoogle(credentials.ClientID, credentials.ClientSecret)

	token, err := oauthClient.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	_, err = c.store.CreateSessionGoogle(ctx, sessionState.UserID, token.AccessToken, token.RefreshToken, sessionState.Scopes, token.Expiry, sessionState.Credentials)
	if err != nil {
		return "", err
	}
	return sessionState.RedirectURL, nil
}

func (c *credentials) RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error {
	if userID == "" {
		return errs.InvalidInput
	}
	return c.store.RemoveSessionGoogle(ctx, userID, credentialsID)
}

func (c *credentials) GetSessionsGoogleByUser(ctx context.Context, userID string) ([]model.SessionGoogle, error) {
	if userID == "" {
		return nil, errs.InvalidInput
	}
	return c.store.GetSessionGoogleAll(ctx, userID)
}
