package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

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
}

var _ Service = &service{}

type service struct {
	store store.Store
}

// New creates a new service.
func New(s store.Store) *service {
	log.Info().Msg("created service")
	return &service{
		store: s,
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
	google, err := s.store.GetCredentialsGoogle(ctx)
	if err != nil {
		return nil, nil, err
	}

	deepL, err := s.store.GetCredentialsDeepL(ctx)
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
