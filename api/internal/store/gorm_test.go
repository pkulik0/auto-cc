package store_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"

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

	credentials, err := s.AddCredentialsDeepL(context.Background(), key)
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

	credentialsID, userID, state, scopes := uint(1), randomString(c), randomString(c), randomString(c)

	err := s.SaveSessionState(context.Background(), credentialsID, userID, state, scopes)
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

	userID, accessToken, refreshToken, expiry, credentialsID, scopes := randomString(c), randomString(c), randomString(c), int64(1), uint(1), randomString(c)

	session, err := s.CreateSessionGoogle(context.Background(), userID, accessToken, refreshToken, expiry, credentialsID, scopes)
	c.Assert(err, qt.IsNil)
	c.Assert(session.UserID, qt.Equals, userID)
	c.Assert(session.AccessToken, qt.Equals, accessToken)
	c.Assert(session.RefreshToken, qt.Equals, refreshToken)
	c.Assert(session.Expiry, qt.Equals, expiry)
	c.Assert(session.CredentialsID, qt.Equals, credentialsID)
	c.Assert(session.Scopes, qt.Equals, scopes)

	sessions, err := s.GetUserSessionsGoogle(context.Background(), userID)
	c.Assert(err, qt.IsNil)
	c.Assert(sessions, qt.Contains, *session)

	err = s.RemoveSessionGoogle(context.Background(), userID, credentialsID)
	c.Assert(err, qt.IsNil)

	sessions, err = s.GetUserSessionsGoogle(context.Background(), userID)
	c.Assert(err, qt.IsNil)
	c.Assert(sessions, qt.Not(qt.Contains), *session)
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
		deepl, err = store.AddCredentialsDeepL(ctx, key)
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
		deepl, err = store.AddCredentialsDeepL(ctx, key)
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
