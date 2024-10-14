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

	db.AutoMigrate(&model.CredentialsGoogle{}, &model.CredentialsDeepL{}, &model.SessionGoogle{}, &model.SessionState{})
	log.Debug().Msg("migrated database models")

	return &gormStore{db: db}, nil
}

func (s *gormStore) Transaction(ctx context.Context, f func(ctx context.Context, store Store) error) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := f(ctx, &gormStore{db: tx})
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

func (s *gormStore) GetCredentialsGoogleAll(ctx context.Context) ([]model.CredentialsGoogle, error) {
	var credentials []model.CredentialsGoogle

	result := s.db.WithContext(ctx).Find(&credentials)
	if result.Error != nil {
		return nil, result.Error
	}

	return credentials, nil
}

func (s *gormStore) GetCredentialsDeepLAll(ctx context.Context) ([]model.CredentialsDeepL, error) {
	var credentials []model.CredentialsDeepL

	result := s.db.WithContext(ctx).Find(&credentials)
	if result.Error != nil {
		return nil, result.Error
	}

	return credentials, nil
}

func (s *gormStore) GetCredentialsGoogleByID(ctx context.Context, id uint) (*model.CredentialsGoogle, error) {
	var credentials model.CredentialsGoogle

	result := s.db.WithContext(ctx).First(&credentials, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &credentials, nil
}

func (s *gormStore) GetCredentialsDeepLByID(ctx context.Context, id uint) (*model.CredentialsDeepL, error) {
	var credentials model.CredentialsDeepL

	result := s.db.WithContext(ctx).First(&credentials, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &credentials, nil
}

func (s *gormStore) RemoveCredentialsGoogle(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.CredentialsGoogle{}, id)
	if result.Error != nil {
		return result.Error
	}

	result = s.db.WithContext(ctx).Where("credentials_id = ?", id).Delete(&model.SessionGoogle{})
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

func (s *gormStore) CreateSessionGoogle(ctx context.Context, userID, accessToken, refreshToken string, expiry int64, credentialsID uint, scopes string) (*model.SessionGoogle, error) {
	session := &model.SessionGoogle{
		UserID:        userID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		Expiry:        expiry,
		CredentialsID: credentialsID,
		Scopes:        scopes,
	}

	result := s.db.WithContext(ctx).Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func (s *gormStore) GetUserSessionsGoogle(ctx context.Context, userID string) ([]model.SessionGoogle, error) {
	var sessions []model.SessionGoogle

	result := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

func (s *gormStore) RemoveSessionGoogle(ctx context.Context, userID string, credentialsID uint) error {
	result := s.db.WithContext(ctx).Where("user_id = ? AND credentials_id = ?", userID, credentialsID).Delete(&model.SessionGoogle{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *gormStore) SaveSessionState(ctx context.Context, credentialsID uint, userID, state string, scopes string) error {
	sessionState := &model.SessionState{
		UserID:        userID,
		State:         state,
		CredentialsID: credentialsID,
		Scopes:        scopes,
	}

	result := s.db.WithContext(ctx).Create(sessionState)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *gormStore) GetSessionState(ctx context.Context, state string) (*model.SessionState, error) {
	var sessionState model.SessionState

	result := s.db.WithContext(ctx).Where("state = ?", state).First(&sessionState)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sessionState, nil
}
