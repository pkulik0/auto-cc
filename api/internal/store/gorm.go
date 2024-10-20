package store

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/model"
	"github.com/pkulik0/autocc/api/internal/quota"
)

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
	log.Debug().Str("host", host).Uint16("port", port).Str("user", user).Str("db", dbName).Msg("connected to psql")

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

func (s *gormStore) AddCredentialsDeepL(ctx context.Context, key string, usage uint) (*model.CredentialsDeepL, error) {
	credentials := &model.CredentialsDeepL{
		Key:   key,
		Usage: usage,
	}

	result := s.db.WithContext(ctx).Create(credentials)
	if result.Error != nil {
		return nil, result.Error
	}

	return credentials, nil
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
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).Delete(&model.CredentialsGoogle{}, id)
		if result.Error != nil {
			return result.Error
		}

		result = tx.WithContext(ctx).Where("credentials_id = ?", id).Delete(&model.SessionGoogle{})
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (s *gormStore) RemoveCredentialsDeepL(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.CredentialsDeepL{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *gormStore) CreateSessionGoogle(ctx context.Context, userID, accessToken, refreshToken, scopes string, expiry time.Time, credentials model.CredentialsGoogle) (*model.SessionGoogle, error) {
	session := &model.SessionGoogle{
		UserID:        userID,
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		Expiry:        expiry,
		CredentialsID: credentials.ID,
		Scopes:        scopes,
	}

	result := s.db.WithContext(ctx).Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	session.Credentials = credentials
	return session, nil
}

func (s *gormStore) GetSessionGoogleByCredentialsID(ctx context.Context, credentialsID uint, userID string) (*model.SessionGoogle, error) {
	var session model.SessionGoogle

	result := s.db.WithContext(ctx).
		Preload("Credentials").
		Where("credentials_id = ? AND user_id = ?", credentialsID, userID).
		First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (s *gormStore) GetSessionGoogleAll(ctx context.Context, userID string) ([]model.SessionGoogle, error) {
	var sessions []model.SessionGoogle

	result := s.db.WithContext(ctx).
		Preload("Credentials").
		Where("user_id = ?", userID).
		Find(&sessions)
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

func (s *gormStore) SaveSessionState(ctx context.Context, credentialsID uint, userID, state, scopes, redirectURL string) error {
	sessionState := &model.SessionState{
		UserID:        userID,
		State:         state,
		CredentialsID: credentialsID,
		Scopes:        scopes,
		RedirectURL:   redirectURL,
	}

	result := s.db.WithContext(ctx).Create(sessionState)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *gormStore) GetSessionState(ctx context.Context, state string) (*model.SessionState, error) {
	var sessionState model.SessionState

	result := s.db.WithContext(ctx).Preload("Credentials").Where("state = ?", state).First(&sessionState)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sessionState, nil
}

func (s *gormStore) GetSessionGoogleByAvailableCost(ctx context.Context, userID string, cost uint) (*model.SessionGoogle, func() error, error) {
	var session model.SessionGoogle

	maxUsage := quota.Google - cost

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Preload("Credentials").
			Joins("JOIN credentials_google ON credentials_google.id = sessions_google.credentials_id").
			Where("user_id = ? AND credentials_google.usage < ?", userID, maxUsage).
			Order("credentials_google.usage").
			Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: "credentials_google"}}).
			First(&session)
		switch result.Error {
		case nil:
		case gorm.ErrRecordNotFound:
			return errs.NotFound
		default:
			return result.Error
		}

		result = tx.Model(&session.Credentials).Update("usage", gorm.Expr("usage + ?", cost))
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	revert := func() error {
		return s.db.WithContext(ctx).Model(&session.Credentials).Update("usage", gorm.Expr("usage - ?", cost)).Error
	}

	session.Credentials.Usage += cost
	return &session, revert, nil
}

func (s *gormStore) UpdateSessionGoogle(ctx context.Context, session *model.SessionGoogle) error {
	result := s.db.WithContext(ctx).Save(session)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *gormStore) GetCredentialsDeepLByAvailableCost(ctx context.Context, cost uint) (*model.CredentialsDeepL, func() error, error) {
	var credentials model.CredentialsDeepL

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("usage < ?", quota.DeepL-cost).
			Order("usage").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&credentials)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&credentials).Update("usage", gorm.Expr("usage + ?", cost))
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	revert := func() error {
		return s.db.WithContext(ctx).Model(&credentials).Update("usage", gorm.Expr("usage - ?", cost)).Error
	}

	credentials.Usage += cost
	return &credentials, revert, nil
}
