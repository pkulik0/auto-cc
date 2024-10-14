package oauth

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// OAuth2Config is an interface for OAuth2 configurations.
type OAuth2Config interface {
	// AuthCodeURL returns a URL for authenticating with OAuth2.
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	// Exchange exchanges an authorization code for a token.
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	// Client returns an HTTP client using the provided token.
	Client(ctx context.Context, t *oauth2.Token) *http.Client
	// TokenSource returns a token source using the provided token.
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

// Configs is an interface for getting Configs configurations.
//
//go:generate mockgen -destination=../mock/oauth.go -package=mock . Configs,OAuth2Config
type Configs interface {
	// GetGoogle returns Google OAuth2 configurations.
	GetGoogle(clientID, clientSecret string) (OAuth2Config, string)
}

var _ Configs = &configs{}

type configs struct {
	googleCallbackURL string
}

// New creates a new Configs client.
func New(googleCallbackURL string) *configs {
	return &configs{
		googleCallbackURL: googleCallbackURL,
	}
}

// GetGoogle returns Google OAuth2 config and scopes.
func (o *configs) GetGoogle(clientID, clientSecret string) (OAuth2Config, string) {
	scopes := []string{youtube.YoutubeForceSslScope}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  o.googleCallbackURL,
	}, strings.Join(scopes, ";")
}

var _ oauth2.TokenSource = &reactiveTokenSource{}

type reactiveTokenSource struct {
	lastToken   *oauth2.Token
	tokenSource oauth2.TokenSource
	onChange    func(*oauth2.Token)
}

// NewReactiveTokenSource creates a new token source that calls onChange when the token changes.
func NewReactiveTokenSource(src oauth2.TokenSource, onChange func(*oauth2.Token)) (*reactiveTokenSource, error) {
	t, err := src.Token()
	if err != nil {
		return nil, err
	}

	return &reactiveTokenSource{
		tokenSource: src,
		lastToken:   t,
		onChange:    onChange,
	}, nil
}

func (s *reactiveTokenSource) Token() (*oauth2.Token, error) {
	t, err := s.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if t.AccessToken != s.lastToken.AccessToken {
		s.onChange(t)
	}
	s.lastToken = t

	return t, nil
}
