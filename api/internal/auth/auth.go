package auth

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
)

// Auth is an interface for authenticating users.
//
//go:generate mockgen -destination=../mock/auth.go -package=mock . Auth
type Auth interface {
	// AuthMiddleware is a middleware that authenticates the user.
	AuthMiddleware(next http.Handler) http.Handler
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
			log.Info().Dur("duration", d).Msg("refreshing token in")

			select {
			case <-timer.C:
				log.Info().Msg("refreshing service account token")

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

	log.Info().Str("realm", realm).Msg("created keycloak auth")
	return &keycloakAuth{
		client: client,
		token:  token,
		realm:  realm,
	}, nil
}

type (
	userIdContextKey struct{}
	userSuperuserKey struct{}
)

func UserFromContext(ctx context.Context) (userId string, isSuperuser bool, ok bool) {
	userId, ok = ctx.Value(userIdContextKey{}).(string)
	if !ok {
		return "", false, false
	}
	isSuperuser, ok = ctx.Value(userSuperuserKey{}).(bool)
	if !ok {
		return "", false, false
	}
	return userId, isSuperuser, true
}

func contextWithUser(ctx context.Context, userId string, isSuperuser bool) context.Context {
	ctx = context.WithValue(ctx, userIdContextKey{}, userId)
	ctx = context.WithValue(ctx, userSuperuserKey{}, isSuperuser)
	return ctx
}

const (
	superusersGroup = "Superusers"
)

func (a *keycloakAuth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" || !strings.HasPrefix(accessToken, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		token, claims, err := a.client.DecodeAccessToken(r.Context(), accessToken, a.realm)
		if err != nil {
			log.Error().Err(err).Msg("failed to decode access token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := claims.GetSubject()
		if err != nil {
			log.Error().Err(err).Msg("failed to get subject")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		groups, err := a.client.GetUserGroups(r.Context(), a.token.AccessToken(), a.realm, userId, gocloak.GetGroupsParams{})
		if err != nil {
			log.Error().Err(err).Msg("failed to get user groups")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		isSuperuser := false
		for _, group := range groups {
			if group.Name == nil {
				continue
			}
			if *group.Name != superusersGroup {
				continue
			}
			isSuperuser = true
			break
		}

		log.Info().Str("user_id", userId).Bool("is_superuser", isSuperuser).Msg("authenticated")
		next.ServeHTTP(w, r.WithContext(contextWithUser(r.Context(), userId, isSuperuser)))
	})
}
