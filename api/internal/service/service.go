package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"github.com/pkulik0/autocc/api/internal/model"
	"github.com/pkulik0/autocc/api/internal/oauth"
	"github.com/pkulik0/autocc/api/internal/store"
)

// Service is an interface for the service layer.
//
//go:generate mockgen -destination=../mock/service.go -package=mock . Service
type Service interface {
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
	GetSessionGoogleURL(ctx context.Context, credentialsID uint, userID string) (string, error)
	// CreateSessionGoogle creates a new Google API session.
	CreateSessionGoogle(ctx context.Context, state, code string) error
	// RemoveSessionGoogle removes a Google API session.
	RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error
	// GetSessionsGoogleByUser returns all Google API sessions for a user.
	GetSessionsGoogleByUser(ctx context.Context, userID string) ([]model.SessionGoogle, error)
}

var _ Service = &service{}

type service struct {
	store store.Store
	oauth oauth.OAuth2
}

// New creates a new service.
func New(s store.Store, o oauth.OAuth2) *service {
	log.Info().Msg("created service")
	return &service{
		store: s,
		oauth: o,
	}
}

var (
	ErrInvalidInput = errors.New("invalid input")
)

func (s *service) AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error) {
	if clientID == "" || clientSecret == "" {
		return nil, ErrInvalidInput
	}

	return s.store.AddCredentialsGoogle(ctx, clientID, clientSecret)
}

func (s *service) AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error) {
	if key == "" {
		return nil, ErrInvalidInput
	}

	return s.store.AddCredentialsDeepL(ctx, key)
}

func (s *service) GetCredentials(ctx context.Context) ([]model.CredentialsGoogle, []model.CredentialsDeepL, error) {
	google, err := s.store.GetCredentialsGoogleAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	deepL, err := s.store.GetCredentialsDeepLAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	return google, deepL, nil
}

func (s *service) RemoveCredentialsGoogle(ctx context.Context, id uint) error {
	return s.store.RemoveCredentialsGoogle(ctx, id)
}

func (s *service) RemoveCredentialsDeepL(ctx context.Context, id uint) error {
	return s.store.RemoveCredentialsDeepL(ctx, id)
}

func generateState() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func (s *service) GetSessionGoogleURL(ctx context.Context, credentialsID uint, userID string) (string, error) {
	if userID == "" {
		return "", ErrInvalidInput
	}

	credentials, err := s.store.GetCredentialsGoogleByID(ctx, credentialsID)
	if err != nil {
		return "", err
	}
	oauthClient, scopes := s.oauth.GetGoogle(credentials.ClientID, credentials.ClientSecret)

	state, err := generateState()
	if err != nil {
		return "", err
	}
	err = s.store.SaveSessionState(ctx, credentialsID, userID, state, scopes)
	if err != nil {
		return "", err
	}

	return oauthClient.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *service) CreateSessionGoogle(ctx context.Context, state, code string) error {
	if state == "" || code == "" {
		return ErrInvalidInput
	}

	sessionState, err := s.store.GetSessionState(ctx, state)
	if err != nil {
		return err
	}

	credentials, err := s.store.GetCredentialsGoogleByID(ctx, sessionState.CredentialsID)
	if err != nil {
		return err
	}
	oauthClient, _ := s.oauth.GetGoogle(credentials.ClientID, credentials.ClientSecret)

	token, err := oauthClient.Exchange(ctx, code)
	if err != nil {
		return err
	}

	_, err = s.store.CreateSessionGoogle(ctx, sessionState.UserID, token.AccessToken, token.RefreshToken, token.Expiry.Unix(), sessionState.CredentialsID, sessionState.Scopes)
	return err
}

func (s *service) RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error {
	if userID == "" {
		return ErrInvalidInput
	}
	return s.store.RemoveSessionGoogle(ctx, userID, credentialsID)
}

func (s *service) GetSessionsGoogleByUser(ctx context.Context, userID string) ([]model.SessionGoogle, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	return s.store.GetUserSessionsGoogle(ctx, userID)
}
