package model

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/pkulik0/autocc/api/internal/oauth"
)

// SessionGoogle is a model for storing Google API sessions.
type SessionGoogle struct {
	gorm.Model
	UserID        string
	AccessToken   string
	RefreshToken  string
	Expiry        time.Time
	CredentialsID uint
	Credentials   CredentialsGoogle `gorm:"foreignKey:CredentialsID"`
	Scopes        string
}

// TableName returns the table name for the model.
func (s *SessionGoogle) TableName() string {
	return "sessions_google"
}

// GetTokenSource returns a token source for the session.
func (s *SessionGoogle) GetTokenSource(ctx context.Context, onChange func(*oauth2.Token)) (oauth2.TokenSource, error) {
	t := &oauth2.Token{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		Expiry:       s.Expiry,
	}

	config, _ := oauth.New("doesnt-matter").GetGoogle(s.Credentials.ClientID, s.Credentials.ClientSecret)

	return oauth.NewReactiveTokenSource(
		oauth2.ReuseTokenSource(t, config.TokenSource(ctx, t)),
		onChange,
	)
}

// SessionState is a model for storing state values used in OAuth2.
type SessionState struct {
	gorm.Model
	UserID        string
	State         string
	CredentialsID uint
	Credentials   CredentialsGoogle `gorm:"foreignKey:CredentialsID"`
	Scopes        string
	RedirectURL   string
}

// TableName returns the table name for the model.
func (s *SessionState) TableName() string {
	return "sessions_state"
}
