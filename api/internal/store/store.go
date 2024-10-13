package store

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pkulik0/autocc/api/internal/model"
)

// Store is an interface for storing external services and client credentials.
//
//go:generate mockgen -destination=mocks/store.go -package=mocks . Store
type Store interface {
	// AddCredentialsGoogle adds Google client credentials to the store.
	AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error)
	// AddCredentialsDeepL adds DeepL client credentials to the store.
	AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error)
	// GetCredentialsGoogle returns all Google client credentials.
	GetCredentialsGoogle(ctx context.Context) ([]model.CredentialsGoogle, error)
	// GetCredentialsDeepL returns all DeepL client credentials.
	GetCredentialsDeepL(ctx context.Context) ([]model.CredentialsDeepL, error)
	// RemoveCredentialsGoogle removes Google client credentials from the store.
	RemoveCredentialsGoogle(ctx context.Context, id uint) error
	// RemoveCredentialsDeepL removes DeepL client credentials from the store.
	RemoveCredentialsDeepL(ctx context.Context, id uint) error
}

var _ Store = &gormStore{}

type gormStore struct {
	db *gorm.DB
}

// New creates a new store with a connection to the PostgreSQL database.
func New(host string, port uint16, user, pass, dbName string) (*gormStore, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Info().Str("host", host).Uint16("port", port).Str("user", user).Str("db", dbName).Msg("connected to psql")

	db.AutoMigrate(&model.CredentialsGoogle{}, &model.CredentialsDeepL{})
	log.Debug().Msg("migrated database models")

	return &gormStore{db: db}, nil
}

func (s *gormStore) AddCredentialsGoogle(ctx context.Context, clientID, clientSecret string) (*model.CredentialsGoogle, error) {
	client := &model.CredentialsGoogle{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	result := s.db.WithContext(ctx).Create(client)
	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}

func (s *gormStore) AddCredentialsDeepL(ctx context.Context, key string) (*model.CredentialsDeepL, error) {
	client := &model.CredentialsDeepL{
		Key: key,
	}

	result := s.db.WithContext(ctx).Create(client)
	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}

func (s *gormStore) GetCredentialsGoogle(ctx context.Context) ([]model.CredentialsGoogle, error) {
	var credentials []model.CredentialsGoogle

	result := s.db.WithContext(ctx).Find(&credentials)
	if result.Error != nil {
		return nil, result.Error
	}

	return credentials, nil
}

func (s *gormStore) GetCredentialsDeepL(ctx context.Context) ([]model.CredentialsDeepL, error) {
	var credentials []model.CredentialsDeepL

	result := s.db.WithContext(ctx).Find(&credentials)
	if result.Error != nil {
		return nil, result.Error
	}

	return credentials, nil
}

func (s *gormStore) RemoveCredentialsGoogle(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.CredentialsGoogle{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *gormStore) RemoveCredentialsDeepL(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.CredentialsDeepL{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
