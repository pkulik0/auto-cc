package credentials_test

import (
	"context"
	"errors"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"

	"github.com/pkulik0/autocc/api/internal/credentials"
	"github.com/pkulik0/autocc/api/internal/mock"
	"github.com/pkulik0/autocc/api/internal/model"
)

func TestService(t *testing.T) {
	c := qt.New(t)

	retErr := errors.New("error")

	testCases := []struct {
		name      string
		setupMock func(mockStore *mock.MockStore, mockOAuth *mock.MockOAuth2Config)
		test      func(c *qt.C, s credentials.Credentials)
	}{
		{
			name: "AddCredentialsGoogle",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				call := store.EXPECT().AddCredentialsGoogle(gomock.Any(), "clientID", "clientSecret").Return(&model.CredentialsGoogle{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				}, nil).Times(1)

				store.EXPECT().AddCredentialsGoogle(gomock.Any(), "clientID", "clientSecret").Return(nil, retErr).Times(1).After(call)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				cred, err := s.AddCredentialsGoogle(context.Background(), "clientID", "clientSecret")
				c.Assert(err, qt.IsNil)
				c.Assert(cred.ClientID, qt.Equals, "clientID")

				_, err = s.AddCredentialsGoogle(context.Background(), "clientID", "clientSecret")
				c.Assert(err, qt.Equals, retErr)

				_, err = s.AddCredentialsGoogle(context.Background(), "", "")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
		{
			name: "AddCredentialsDeepL",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				call := store.EXPECT().AddCredentialsDeepL(gomock.Any(), "key").Return(&model.CredentialsDeepL{
					Key: "key",
				}, nil).Times(1)

				store.EXPECT().AddCredentialsDeepL(gomock.Any(), "key").Return(nil, retErr).Times(1).After(call)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				cred, err := s.AddCredentialsDeepL(context.Background(), "key")
				c.Assert(err, qt.IsNil)
				c.Assert(cred.Key, qt.Equals, "key")

				_, err = s.AddCredentialsDeepL(context.Background(), "key")
				c.Assert(err, qt.Equals, retErr)

				_, err = s.AddCredentialsDeepL(context.Background(), "")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
		{
			name: "GetCredentials",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				store.EXPECT().GetCredentialsGoogleAll(gomock.Any()).Return([]model.CredentialsGoogle{{ClientID: "clientID", ClientSecret: "clientSecret"}}, nil).Times(1)
				store.EXPECT().GetCredentialsDeepLAll(gomock.Any()).Return([]model.CredentialsDeepL{{Key: "key1"}}, nil).Times(1)

				store.EXPECT().GetCredentialsGoogleAll(gomock.Any()).Return(nil, retErr).Times(1)

				store.EXPECT().GetCredentialsGoogleAll(gomock.Any()).Return([]model.CredentialsGoogle{{ClientID: "clientID", ClientSecret: "clientSecret"}}, nil).Times(1)
				store.EXPECT().GetCredentialsDeepLAll(gomock.Any()).Return(nil, retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				g, d, err := s.GetCredentials(context.Background())
				c.Assert(err, qt.IsNil)
				c.Assert(g, qt.DeepEquals, []model.CredentialsGoogle{{ClientID: "clientID", ClientSecret: "clientSecret"}})
				c.Assert(d, qt.DeepEquals, []model.CredentialsDeepL{{Key: "key1"}})

				_, _, err = s.GetCredentials(context.Background())
				c.Assert(err, qt.Equals, retErr)

				_, _, err = s.GetCredentials(context.Background())
				c.Assert(err, qt.Equals, retErr)
			},
		},
		{
			name: "RemoveCredentialsGoogle",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				store.EXPECT().RemoveCredentialsGoogle(gomock.Any(), uint(1)).Return(nil).Times(1)
				store.EXPECT().RemoveCredentialsGoogle(gomock.Any(), uint(1)).Return(retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				err := s.RemoveCredentialsGoogle(context.Background(), 1)
				c.Assert(err, qt.IsNil)

				err = s.RemoveCredentialsGoogle(context.Background(), 1)
				c.Assert(err, qt.Equals, retErr)
			},
		},
		{
			name: "RemoveCredentialsDeepL",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				store.EXPECT().RemoveCredentialsDeepL(gomock.Any(), uint(1)).Return(nil).Times(1)
				store.EXPECT().RemoveCredentialsDeepL(gomock.Any(), uint(1)).Return(retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				err := s.RemoveCredentialsDeepL(context.Background(), 1)
				c.Assert(err, qt.IsNil)

				err = s.RemoveCredentialsDeepL(context.Background(), 1)
				c.Assert(err, qt.Equals, retErr)
			},
		},
		{
			name: "GetSessionGoogleURL",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				oauth.EXPECT().AuthCodeURL(gomock.Any(), oauth2.AccessTypeOffline).Return("url").Times(1)
				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(&model.CredentialsGoogle{ClientID: "clientID"}, nil).Times(1)
				store.EXPECT().SaveSessionState(gomock.Any(), uint(1), "userID", gomock.Any(), gomock.Any(), "http://example.com").Return(nil).Times(1)

				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(nil, retErr).Times(1)

				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(&model.CredentialsGoogle{ClientID: "clientID"}, nil).Times(1)
				store.EXPECT().SaveSessionState(gomock.Any(), uint(1), "userID", gomock.Any(), gomock.Any(), "http://example.com").Return(retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				_, err := s.GetSessionGoogleURL(context.Background(), 1, "userID", "http://example.com")
				c.Assert(err, qt.IsNil)

				_, err = s.GetSessionGoogleURL(context.Background(), 1, "userID", "http://example.com")
				c.Assert(err, qt.IsNotNil)

				_, err = s.GetSessionGoogleURL(context.Background(), 1, "userID", "http://example.com")
				c.Assert(err, qt.Equals, retErr)

				_, err = s.GetSessionGoogleURL(context.Background(), 1, "", "http://example.com")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)

				_, err = s.GetSessionGoogleURL(context.Background(), 1, "abc", "invalid_url")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
		{
			name: "CreateSessionGoogle",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				sessState := &model.SessionState{
					CredentialsID: 1,
					UserID:        "userID",
					Scopes:        "scopes",
					RedirectURL:   "http://example.com",
				}
				cred := &model.CredentialsGoogle{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				}
				token := &oauth2.Token{
					AccessToken:  "access",
					RefreshToken: "refresh",
					Expiry:       time.Now(),
				}

				// 1
				store.EXPECT().GetSessionState(gomock.Any(), "state").Return(sessState, nil).Times(1)
				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(cred, nil).Times(1)
				oauth.EXPECT().Exchange(gomock.Any(), "code").Return(token, nil).Times(1)
				store.EXPECT().CreateSessionGoogle(gomock.Any(), "userID", "access", "refresh", "scopes", gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				// 2
				store.EXPECT().GetSessionState(gomock.Any(), "state").Return(nil, retErr).Times(1)

				// 3
				store.EXPECT().GetSessionState(gomock.Any(), "state").Return(sessState, nil).Times(1)
				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(nil, retErr).Times(1)

				// 4
				store.EXPECT().GetSessionState(gomock.Any(), "state").Return(sessState, nil).Times(1)
				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(cred, nil).Times(1)
				oauth.EXPECT().Exchange(gomock.Any(), "code").Return(nil, retErr).Times(1)

				// 5
				store.EXPECT().GetSessionState(gomock.Any(), "state").Return(sessState, nil).Times(1)
				store.EXPECT().GetCredentialsGoogleByID(gomock.Any(), uint(1)).Return(cred, nil).Times(1)
				oauth.EXPECT().Exchange(gomock.Any(), "code").Return(token, nil).Times(1)
				store.EXPECT().CreateSessionGoogle(gomock.Any(), "userID", "access", "refresh", "scopes", gomock.Any(), gomock.Any()).Return(nil, retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				url, err := s.CreateSessionGoogle(context.Background(), "state", "code")
				c.Assert(err, qt.IsNil)
				c.Assert(url, qt.Equals, "http://example.com")

				url, err = s.CreateSessionGoogle(context.Background(), "state", "code")
				c.Assert(err, qt.Equals, retErr)

				url, err = s.CreateSessionGoogle(context.Background(), "state", "code")
				c.Assert(err, qt.Equals, retErr)

				url, err = s.CreateSessionGoogle(context.Background(), "state", "code")
				c.Assert(err, qt.Equals, retErr)

				url, err = s.CreateSessionGoogle(context.Background(), "state", "code")
				c.Assert(err, qt.Equals, retErr)

				url, err = s.CreateSessionGoogle(context.Background(), "", "")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
		{
			name: "RemoveSessionGoogle",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				store.EXPECT().RemoveSessionGoogle(gomock.Any(), "userID", uint(1)).Return(nil).Times(1)
				store.EXPECT().RemoveSessionGoogle(gomock.Any(), "userID", uint(1)).Return(retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				err := s.RemoveSessionGoogle(context.Background(), "userID", 1)
				c.Assert(err, qt.IsNil)

				err = s.RemoveSessionGoogle(context.Background(), "userID", 1)
				c.Assert(err, qt.Equals, retErr)

				err = s.RemoveSessionGoogle(context.Background(), "", 1)
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
		{
			name: "GetSessionsGoogleByUser",
			setupMock: func(store *mock.MockStore, oauth *mock.MockOAuth2Config) {
				store.EXPECT().GetSessionGoogleAll(gomock.Any(), "userID").Return([]model.SessionGoogle{{CredentialsID: 1}}, nil).Times(1)
				store.EXPECT().GetSessionGoogleAll(gomock.Any(), "userID").Return(nil, retErr).Times(1)
			},
			test: func(c *qt.C, s credentials.Credentials) {
				sessions, err := s.GetSessionsGoogleByUser(context.Background(), "userID")
				c.Assert(err, qt.IsNil)
				c.Assert(sessions, qt.DeepEquals, []model.SessionGoogle{{CredentialsID: 1}})

				_, err = s.GetSessionsGoogleByUser(context.Background(), "userID")
				c.Assert(err, qt.Equals, retErr)

				_, err = s.GetSessionsGoogleByUser(context.Background(), "")
				c.Assert(err, qt.Equals, credentials.ErrInvalidInput)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)
			store := mock.NewMockStore(ctrl)
			oauth := mock.NewMockConfigs(ctrl)
			oauthGoogle := mock.NewMockOAuth2Config(ctrl)
			oauth.EXPECT().GetGoogle(gomock.Any(), gomock.Any()).Return(oauthGoogle, "scopes").AnyTimes()
			tc.setupMock(store, oauthGoogle)

			s := credentials.New(store, oauth)
			tc.test(c, s)
		})
	}
}
