package oauth

import (
	"context"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// OAuth2Client is an interface for OAuth2 configurations.
type OAuth2Client interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

// OAuth2 is an interface for getting OAuth2 configurations.
//
//go:generate mockgen -destination=../mock/oauth.go -package=mock . OAuth2,OAuth2Client
type OAuth2 interface {
	GetGoogle(clientID, clientSecret string) (OAuth2Client, string)
}

var _ OAuth2 = &oauth2Client{}

type oauth2Client struct {
	googleCallbackURL string
}

// New creates a new OAuth2 client.
func New(googleCallbackURL string) *oauth2Client {
	return &oauth2Client{
		googleCallbackURL: googleCallbackURL,
	}
}

func (o *oauth2Client) GetGoogle(clientID, clientSecret string) (OAuth2Client, string) {
	scopes := []string{youtube.YoutubeForceSslScope}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  o.googleCallbackURL,
	}, strings.Join(scopes, ";")
}
