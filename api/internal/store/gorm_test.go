package store_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/pkulik0/autocc/api/internal/model"
	"github.com/pkulik0/autocc/api/internal/store"
)

func randomString(c *qt.C) string {
	s := make([]byte, 16)
	_, err := rand.Read(s)
	c.Assert(err, qt.IsNil)
	return hex.EncodeToString(s)
}

func setupStore(c *qt.C) store.Store {
	s, err := store.New("localhost", 5432, "autocc", "autocc", "autocc")
	c.Assert(err, qt.IsNil)
	return s
}

func TestConnection(t *testing.T) {
	c := qt.New(t)
	_ = setupStore(c)

	_, err := store.New("doesnt-exist", 5432, "autocc", "autocc", "autocc")
	c.Assert(err, qt.IsNotNil)
}

func TestCredentialsGoogle(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	id, secret := randomString(c), randomString(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), id, secret)
	c.Assert(err, qt.IsNil)
	c.Assert(credentials.ClientID, qt.Equals, id)
	c.Assert(credentials.ClientSecret, qt.Equals, secret)

	retrieved, err := s.GetCredentialsGoogleByID(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved, qt.DeepEquals, credentials)

	credentialsAll, err := s.GetCredentialsGoogleAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(credentialsAll, qt.Contains, *credentials)

	err = s.RemoveCredentialsGoogle(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNil)

	retrieved, err = s.GetCredentialsGoogleByID(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNotNil)
}

func TestCredentialsDeepL(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	key := randomString(c)

	credentials, err := s.AddCredentialsDeepL(context.Background(), key, 5)
	c.Assert(err, qt.IsNil)
	c.Assert(credentials.Key, qt.Equals, key)

	retrieved, err := s.GetCredentialsDeepLByID(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved, qt.DeepEquals, credentials)

	credentialsAll, err := s.GetCredentialsDeepLAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(credentialsAll, qt.Contains, *credentials)

	err = s.RemoveCredentialsDeepL(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNil)

	retrieved, err = s.GetCredentialsDeepLByID(context.Background(), credentials.ID)
	c.Assert(err, qt.IsNotNil)
}

func TestSessionState(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentialsID, userID, state, scopes, url := uint(1), randomString(c), randomString(c), randomString(c), "https://example.com"

	err := s.SaveSessionState(context.Background(), credentialsID, userID, state, scopes, url)
	c.Assert(err, qt.IsNil)

	retrieved, err := s.GetSessionState(context.Background(), state)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved.CredentialsID, qt.Equals, credentialsID)
	c.Assert(retrieved.UserID, qt.Equals, userID)
	c.Assert(retrieved.State, qt.Equals, state)
	c.Assert(retrieved.Scopes, qt.Equals, scopes)

	_, err = s.GetSessionState(context.Background(), "invalid")
	c.Assert(err, qt.IsNotNil)
}

func TestSessionGoogle(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), "client", "secret")
	c.Assert(err, qt.IsNil)

	userID, accessToken, refreshToken, expiry, scopes := randomString(c), randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)

	session, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)
	c.Assert(session.UserID, qt.Equals, userID)
	c.Assert(session.AccessToken, qt.Equals, accessToken)
	c.Assert(session.RefreshToken, qt.Equals, refreshToken)
	c.Assert(session.Expiry, qt.Equals, expiry)
	c.Assert(session.CredentialsID, qt.Equals, credentials.ID)
	c.Assert(session.Scopes, qt.Equals, scopes)

	retrieved, err := s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved, qt.DeepEquals, session)

	err = s.RemoveSessionGoogle(context.Background(), userID, credentials.ID)
	c.Assert(err, qt.IsNil)

	_, err = s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNotNil)
}

func TestUpdateSessionGoogle(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), "client", "secret")
	c.Assert(err, qt.IsNil)

	userID, accessToken, refreshToken, expiry, scopes := randomString(c), randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)

	session, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)

	newAccessToken, newRefreshToken := randomString(c), randomString(c)
	session.AccessToken = newAccessToken
	session.RefreshToken = newRefreshToken

	err = s.UpdateSessionGoogle(context.Background(), session)
	c.Assert(err, qt.IsNil)

	retrieved, err := s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved, qt.DeepEquals, session)
}

func TestGetSessionGoogleAll(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), "client", "secret")
	c.Assert(err, qt.IsNil)

	userID, accessToken, refreshToken, expiry, scopes := randomString(c), randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)
	session1, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)

	accessToken, refreshToken, expiry, scopes = randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)
	session2, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)

	sessions, err := s.GetSessionGoogleAll(context.Background(), userID)
	c.Assert(err, qt.IsNil)

	c.Assert(sessions, qt.Contains, *session1)
	c.Assert(sessions, qt.Contains, *session2)
}

func TestGetSessionGoogleByCredentialsID(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), "client", "secret")
	c.Assert(err, qt.IsNil)

	userID, accessToken, refreshToken, expiry, scopes := randomString(c), randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)
	session, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)

	retrieved, err := s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved, qt.DeepEquals, session)
}

func TestGetSessionGoogleByAvailableCost(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	credentials, err := s.AddCredentialsGoogle(context.Background(), "client", "secret")
	c.Assert(err, qt.IsNil)

	userID, accessToken, refreshToken, expiry, scopes := randomString(c), randomString(c), randomString(c), time.Now().Round(time.Minute), randomString(c)

	session, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, scopes, expiry, *credentials)
	c.Assert(err, qt.IsNil)

	newSession, revert, err := s.GetSessionGoogleByAvailableCost(context.Background(), userID, 1000)
	c.Assert(err, qt.IsNil)
	c.Assert(newSession.ID, qt.DeepEquals, session.ID)
	c.Assert(newSession.Credentials.Usage, qt.Equals, uint(1000))

	retrieved, err := s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved.Credentials.Usage, qt.Equals, uint(1000))

	retrieved, _, err = s.GetSessionGoogleByAvailableCost(context.Background(), userID, 500)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved.ID, qt.DeepEquals, session.ID)
	c.Assert(retrieved.Credentials.Usage, qt.Equals, uint(1500))

	err = revert()
	c.Assert(err, qt.IsNil)

	retrieved, err = s.GetSessionGoogleByCredentialsID(context.Background(), credentials.ID, userID)
	c.Assert(err, qt.IsNil)
	c.Assert(retrieved.Credentials.Usage, qt.Equals, uint(500))
}

func TestTransaction(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	var google *model.CredentialsGoogle
	var deepl *model.CredentialsDeepL

	err := s.Transaction(context.Background(), func(ctx context.Context, store store.Store) error {
		id, secret := randomString(c), randomString(c)
		var err error
		google, err = store.AddCredentialsGoogle(ctx, id, secret)
		if err != nil {
			return err
		}

		key := randomString(c)
		deepl, err = store.AddCredentialsDeepL(ctx, key, 10)
		return err
	})
	c.Assert(err, qt.IsNil)

	googleAll, err := s.GetCredentialsGoogleAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(googleAll, qt.Contains, *google)

	deeplAll, err := s.GetCredentialsDeepLAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(deeplAll, qt.Contains, *deepl)

	err = s.Transaction(context.Background(), func(ctx context.Context, store store.Store) error {
		id, secret := randomString(c), randomString(c)
		var err error
		google, err = store.AddCredentialsGoogle(ctx, id, secret)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		key := randomString(c)
		deepl, err = store.AddCredentialsDeepL(ctx, key, 20)
		return err
	})
	c.Assert(err, qt.IsNotNil)

	newGoogleAll, err := s.GetCredentialsGoogleAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(googleAll, qt.DeepEquals, newGoogleAll)

	newDeepLAll, err := s.GetCredentialsDeepLAll(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(deeplAll, qt.DeepEquals, newDeepLAll)
}

func TestContext(t *testing.T) {
	c := qt.New(t)
	s := setupStore(c)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.AddCredentialsGoogle(ctx, randomString(c), randomString(c))
	c.Assert(err, qt.IsNotNil)

	_, err = s.AddCredentialsDeepL(ctx, randomString(c), 0)
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetCredentialsGoogleAll(ctx)
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetCredentialsDeepLAll(ctx)
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetCredentialsGoogleByID(ctx, 1)
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetCredentialsDeepLByID(ctx, 1)
	c.Assert(err, qt.IsNotNil)

	err = s.RemoveCredentialsGoogle(ctx, 1)
	c.Assert(err, qt.IsNotNil)

	err = s.RemoveCredentialsDeepL(ctx, 1)
	c.Assert(err, qt.IsNotNil)

	_, err = s.CreateSessionGoogle(ctx, randomString(c), randomString(c), randomString(c), randomString(c), time.Now(), model.CredentialsGoogle{})
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetSessionGoogleAll(ctx, randomString(c))
	c.Assert(err, qt.IsNotNil)

	err = s.RemoveSessionGoogle(ctx, randomString(c), uint(1))
	c.Assert(err, qt.IsNotNil)

	err = s.SaveSessionState(ctx, uint(1), randomString(c), randomString(c), randomString(c), "https://example.com")
	c.Assert(err, qt.IsNotNil)

	_, err = s.GetSessionState(ctx, randomString(c))
	c.Assert(err, qt.IsNotNil)

	_, _, err = s.GetSessionGoogleByAvailableCost(ctx, randomString(c), 1)
	c.Assert(err, qt.IsNotNil)

	_ = s.UpdateSessionGoogle(ctx, &model.SessionGoogle{})
	c.Assert(err, qt.IsNotNil)

	err = s.Transaction(ctx, func(ctx context.Context, store store.Store) error {
		return nil
	})
	c.Assert(err, qt.IsNotNil)
}
