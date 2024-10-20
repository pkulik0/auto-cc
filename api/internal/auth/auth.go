package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
)

// Auth is an interface for authenticating users.
//
//go:generate mockgen -destination=../mock/auth.go -package=mock . Auth
type Auth interface {
	// Authenticate authenticates based on the access token.
	Authenticate(ctx context.Context, accessToken string) (userID string, isSuperuser bool, err error)
}

var _ Auth = &keycloakAuth{}

type token struct {
	token  *gocloak.JWT
	mutex  sync.Mutex
	stopCh chan struct{}
}

func newToken(client *gocloak.GoCloak, realm, clientID, clientSecret string) (*token, error) {
	fetch := func() (*gocloak.JWT, error) {
		return client.LoginClient(context.Background(), clientID, clientSecret, realm)
	}
	clientToken, err := fetch()
	if err != nil {
		return nil, err
	}

	t := &token{
		token:  clientToken,
		stopCh: make(chan struct{}),
	}

	go func() {
		getDuration := func() time.Duration {
			return time.Duration(t.token.ExpiresIn-10) * time.Second
		}
		timer := time.NewTimer(time.Hour)

		for {
			d := getDuration()
			timer.Reset(d)
			log.Debug().Dur("duration", d).Msg("refreshing token in")

			select {
			case <-timer.C:
				log.Debug().Msg("refreshing service account token")

				func() {
					t.mutex.Lock()
					defer t.mutex.Unlock()

					token, err := fetch()
					if err != nil {
						log.Error().Err(err).Msg("failed to fetch token")
						return
					}
					t.token = token
				}()
			case <-t.stopCh:
				return
			}
		}
	}()

	return t, nil
}

func (t *token) AccessToken() string {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.token.AccessToken
}

type keycloakAuth struct {
	client *gocloak.GoCloak
	token  *token
	realm  string
}

// New creates a new instance of keycloak auth.
func New(ctx context.Context, url, realm, clientID, clientSecret string) (*keycloakAuth, error) {
	client := gocloak.NewClient(url)

	token, err := newToken(client, realm, clientID, clientSecret)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("realm", realm).Msg("created keycloak auth")
	return &keycloakAuth{
		client: client,
		token:  token,
		realm:  realm,
	}, nil
}

type (
	contextKeyUserID    struct{}
	contextKeySuperuser struct{}
)

func UserFromContext(ctx context.Context) (userId string, isSuperuser bool, ok bool) {
	userId, ok = ctx.Value(contextKeyUserID{}).(string)
	if !ok {
		return "", false, false
	}
	isSuperuser, ok = ctx.Value(contextKeySuperuser{}).(bool)
	if !ok {
		return "", false, false
	}
	return userId, isSuperuser, true
}

func ContextWithUser(ctx context.Context, userId string, isSuperuser bool) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID{}, userId)
	ctx = context.WithValue(ctx, contextKeySuperuser{}, isSuperuser)
	return ctx
}

const (
	superusersGroup = "Superusers"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func (a *keycloakAuth) Authenticate(ctx context.Context, accessToken string) (string, bool, error) {
	token, claims, err := a.client.DecodeAccessToken(ctx, accessToken, a.realm)
	if err != nil {
		return "", false, err
	}
	if !token.Valid {
		return "", false, ErrInvalidToken
	}
	userId, err := claims.GetSubject()
	if err != nil {
		return "", false, ErrInvalidToken
	}

	groups, err := a.client.GetUserGroups(ctx, a.token.AccessToken(), a.realm, userId, gocloak.GetGroupsParams{})
	if err != nil {
		return "", false, err
	}

	isSuperuser := false
	for _, group := range groups {
		if group.Name != nil && *group.Name == superusersGroup {
			isSuperuser = true
			break
		}
	}

	log.Debug().Str("user_id", userId).Bool("is_superuser", isSuperuser).Msg("authenticated")
	return userId, isSuperuser, nil
}
