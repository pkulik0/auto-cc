package model

import "gorm.io/gorm"

// SessionGoogle is a model for storing Google API sessions.
type SessionGoogle struct {
	gorm.Model
	UserID        string
	AccessToken   string
	RefreshToken  string
	Expiry        int64
	CredentialsID uint
	Credentials   CredentialsGoogle `gorm:"foreignKey:CredentialsID"`
	Scopes        string
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
