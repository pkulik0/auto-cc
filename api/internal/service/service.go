package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"

	"github.com/pkulik0/autocc/api/internal/model"
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
	store             store.Store
	googleCallbackURL string
}

// New creates a new service.
func New(s store.Store, googleCallbackURL string) *service {
	log.Info().Msg("created service")
	return &service{
		store:             s,
		googleCallbackURL: googleCallbackURL,
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

func (s *service) getGoogleOauthConfig(clientID, clientSecret string) (*oauth2.Config, string) {
	scopes := []string{youtube.YoutubeForceSslScope}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  s.googleCallbackURL,
	}, strings.Join(scopes, ";")
}

func (s *service) GetSessionGoogleURL(ctx context.Context, credentialsID uint, userID string) (string, error) {
	credentials, err := s.store.GetCredentialsGoogleByID(ctx, credentialsID)
	if err != nil {
		return "", err
	}
	oauthConfig, scopes := s.getGoogleOauthConfig(credentials.ClientID, credentials.ClientSecret)

	state, err := generateState()
	if err != nil {
		return "", err
	}
	err = s.store.SaveSessionState(ctx, credentialsID, userID, state, scopes)
	if err != nil {
		return "", err
	}

	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *service) CreateSessionGoogle(ctx context.Context, state, code string) error {
	sessionState, err := s.store.GetSessionState(ctx, state)
	if err != nil {
		return err
	}

	credentials, err := s.store.GetCredentialsGoogleByID(ctx, sessionState.CredentialsID)
	if err != nil {
		return err
	}
	oauthConfig, _ := s.getGoogleOauthConfig(credentials.ClientID, credentials.ClientSecret)

	token, err := oauthConfig.Exchange(ctx, code)
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
