package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/credentials"
	"github.com/pkulik0/autocc/api/internal/mock"
	"github.com/pkulik0/autocc/api/internal/model"
	"github.com/pkulik0/autocc/api/internal/pb"
)

func TestHandlerRoot(t *testing.T) {
	c := qt.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	s := New(nil, nil, nil, nil)
	s.handlerRoot(w, r)

	c.Assert(w.Code, qt.Equals, http.StatusOK)
	c.Assert(w.Body.String(), qt.Contains, "version")
	c.Assert(w.Body.String(), qt.Contains, "build_time")
}

func TestHandlerCredentials(t *testing.T) {
	c := qt.New(t)

	deepl := []model.CredentialsDeepL{
		{Key: "key1111111111", Usage: 1},
		{Key: "key2222222222", Usage: 2},
	}
	google := []model.CredentialsGoogle{
		{ClientID: "id1", ClientSecret: "secret1111111", Usage: 1},
		{ClientID: "id2", ClientSecret: "secret9999999", Usage: 2},
	}
	retErr := errors.New("error")

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetCredentials(gomock.Any()).Return(google, deepl, nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/credentials", nil)

				server.handlerCredentials(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
				var resp pb.GetCredentialsResponse
				err := proto.Unmarshal(w.Body.Bytes(), &resp)
				c.Assert(err, qt.IsNil)

				c.Assert(resp.Google, qt.HasLen, 2)
				for i, g := range google {
					c.Assert(resp.Google[i].ClientId, qt.Equals, g.ClientID)
					idxOfStar := strings.Index(resp.Google[i].ClientSecret, "*")
					c.Assert(resp.Google[i].ClientSecret[:idxOfStar], qt.Equals, g.ClientSecret[:idxOfStar])
				}

				c.Assert(resp.Deepl, qt.HasLen, 2)
				for i, d := range deepl {
					idxOfStar := strings.Index(resp.Deepl[i].Key, "*")
					c.Assert(resp.Deepl[i].Key[:idxOfStar], qt.Equals, d.Key[:idxOfStar])
				}
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetCredentials(gomock.Any()).Return(nil, nil, retErr)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/credentials", nil)

				server.handlerCredentials(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestAddCredentialsGoogle(t *testing.T) {
	c := qt.New(t)

	req := &pb.AddCredentialsGoogleRequest{
		ClientId:     "id",
		ClientSecret: "secret",
	}
	data, err := proto.Marshal(req)
	c.Assert(err, qt.IsNil)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsGoogle(gomock.Any(), "id", "secret").Return(&model.CredentialsGoogle{ClientID: "id", ClientSecret: "secret"}, nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/google", bytes.NewReader(data))

				server.handlerAddCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
				var resp pb.AddCredentialsGoogleResponse
				err := proto.Unmarshal(w.Body.Bytes(), &resp)
				c.Assert(err, qt.IsNil)

				c.Assert(resp.Credentials.ClientId, qt.Equals, "id")
				idxOfStar := strings.Index(resp.Credentials.ClientSecret, "*")
				c.Assert(resp.Credentials.ClientSecret[:idxOfStar], qt.Equals, "secret"[:idxOfStar])
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsGoogle(gomock.Any(), "id", "secret").Return(nil, errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/google", bytes.NewReader(data))

				server.handlerAddCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().AddCredentialsGoogle(gomock.Any(), "id", "secret").Return(nil, credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/google", bytes.NewReader(data))

				server.handlerAddCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "invalid request",
			setupMocks: func(s *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/google", strings.NewReader("invalid"))

				server.handlerAddCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestAddCredentialsDeepL(t *testing.T) {
	c := qt.New(t)

	key := "key123123123"
	req := &pb.AddCredentialsDeepLRequest{
		Key: key,
	}
	data, err := proto.Marshal(req)
	c.Assert(err, qt.IsNil)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsDeepL(gomock.Any(), key).Return(&model.CredentialsDeepL{Key: key}, nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/deepl", bytes.NewReader(data))

				server.handlerAddCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
				var resp pb.AddCredentialsDeepLResponse
				err := proto.Unmarshal(w.Body.Bytes(), &resp)
				c.Assert(err, qt.IsNil)

				idxOfStar := strings.Index(resp.Credentials.Key, "*")
				c.Assert(resp.Credentials.Key[:idxOfStar], qt.Equals, key[:idxOfStar])
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsDeepL(gomock.Any(), key).Return(nil, errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/deepl", bytes.NewReader(data))

				server.handlerAddCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().AddCredentialsDeepL(gomock.Any(), key).Return(nil, credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/deepl", bytes.NewReader(data))

				server.handlerAddCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "invalid request",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/credentials/deepl", strings.NewReader("invalid"))

				server.handlerAddCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestRemoveCredentialsGoogle(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsGoogle(gomock.Any(), uint(1)).Return(nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/google/1", nil)
				r.SetPathValue("id", "1")

				server.handlerRemoveCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusNoContent)
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsGoogle(gomock.Any(), uint(1)).Return(errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/google/1", nil)
				r.SetPathValue("id", "1")

				server.handlerRemoveCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name:       "invalid id",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/google/invalid", nil)
				r.SetPathValue("id", "invalid")

				server.handlerRemoveCredentialsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestRemoveCredentialsDeepL(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsDeepL(gomock.Any(), uint(1)).Return(nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/deepl/1", nil)
				r.SetPathValue("id", "1")

				server.handlerRemoveCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusNoContent)
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsDeepL(gomock.Any(), uint(1)).Return(errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/deepl/1", nil)
				r.SetPathValue("id", "1")

				server.handlerRemoveCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name:       "invalid id",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/credentials/deepl/invalid", nil)
				r.SetPathValue("id", "invalid")

				server.handlerRemoveCredentialsDeepL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestHandlerSessionGoogleURL(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionGoogleURL(gomock.Any(), uint(1), "userID", "redirectURL").Return("url", nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/1?redirect_url=redirectURL", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerSessionGoogleURL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
				var resp pb.GetSessionGoogleURLResponse
				err := proto.Unmarshal(w.Body.Bytes(), &resp)
				c.Assert(err, qt.IsNil)

				c.Assert(resp.Url, qt.Equals, "url")
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionGoogleURL(gomock.Any(), uint(1), "userID", "redirectURL").Return("", errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/1?redirect_url=redirectURL", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerSessionGoogleURL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().GetSessionGoogleURL(gomock.Any(), uint(1), "userID", "redirectURL").Return("", credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/1?redirect_url=redirectURL", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerSessionGoogleURL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "invalid id",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/invalid?redirect_url=redirectURL", nil)
				r.SetPathValue("id", "invalid")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerSessionGoogleURL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "no user",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/1?redirect_url=redirectURL", nil)
				r.SetPathValue("id", "1")

				server.handlerSessionGoogleURL(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestHandlerSessionGoogleCallback(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().CreateSessionGoogle(gomock.Any(), "state", "code").Return("http://example.com", nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/callback", nil)
				r.Form = map[string][]string{
					"state": {"state"},
					"code":  {"code"},
				}

				server.handlerSessionGoogleCallback(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusFound)
				c.Assert(w.Header().Get("Location"), qt.Equals, "http://example.com")
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().CreateSessionGoogle(gomock.Any(), "state", "code").Return("", errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/callback", nil)
				r.Form = map[string][]string{
					"state": {"state"},
					"code":  {"code"},
				}

				server.handlerSessionGoogleCallback(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().CreateSessionGoogle(gomock.Any(), "", "").Return("", credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/sessions/google/callback", nil)
				r.Form = map[string][]string{
					"state": {""},
					"code":  {""},
				}

				server.handlerSessionGoogleCallback(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestHandlerUserSessionsGoogle(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionsGoogleByUser(gomock.Any(), "userID").Return([]model.SessionGoogle{{CredentialsID: 123}}, nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/sessions/google", nil)
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerUserSessionsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
				var resp pb.GetUserSessionsGoogleResponse
				err := proto.Unmarshal(w.Body.Bytes(), &resp)
				c.Assert(err, qt.IsNil)

				c.Assert(resp.CredentialIds, qt.HasLen, 1)
				c.Assert(resp.CredentialIds[0], qt.Equals, uint64(123))
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionsGoogleByUser(gomock.Any(), "userID").Return(nil, errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/sessions/google", nil)
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerUserSessionsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().GetSessionsGoogleByUser(gomock.Any(), "userID").Return(nil, credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/sessions/google", nil)
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerUserSessionsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "no user",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/sessions/google", nil)

				server.handlerUserSessionsGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestHandlerRemoveSessionGoogle(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name       string
		setupMocks func(service *mock.MockCredentials)
		test       func(c *qt.C, s *server)
	}{
		{
			name: "success",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveSessionGoogle(gomock.Any(), "userID", uint(1)).Return(nil)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/sessions/google/1", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerRemoveSessionGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusNoContent)
			},
		},
		{
			name: "error",
			setupMocks: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveSessionGoogle(gomock.Any(), "userID", uint(1)).Return(errors.New("error"))
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/sessions/google/1", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerRemoveSessionGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
		{
			name: "invalid input",
			setupMocks: func(s *mock.MockCredentials) {
				s.EXPECT().RemoveSessionGoogle(gomock.Any(), "userID", uint(1)).Return(credentials.ErrInvalidInput)
			},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/sessions/google/1", nil)
				r.SetPathValue("id", "1")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerRemoveSessionGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "invalid id",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/sessions/google/invalid", nil)
				r.SetPathValue("id", "invalid")
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				server.handlerRemoveSessionGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusBadRequest)
			},
		},
		{
			name:       "no user",
			setupMocks: func(service *mock.MockCredentials) {},
			test: func(c *qt.C, server *server) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("DELETE", "/sessions/google/1", nil)
				r.SetPathValue("id", "1")

				server.handlerRemoveSessionGoogle(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(c)

			service := mock.NewMockCredentials(ctrl)
			tc.setupMocks(service)

			s := New(service, nil, nil, nil)
			tc.test(c, s)
		})
	}
}

func TestSuperuserMiddleware(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name string
		test func(c *qt.C)
	}{
		{
			name: "superuser",
			test: func(c *qt.C) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/", nil)
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", true))

				superuserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})).ServeHTTP(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusOK)
			},
		},
		{
			name: "not superuser",
			test: func(c *qt.C) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/", nil)
				r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", false))

				superuserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})).ServeHTTP(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusForbidden)
			},
		},
		{
			name: "no user",
			test: func(c *qt.C) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/", nil)

				superuserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})).ServeHTTP(w, r)

				c.Assert(w.Code, qt.Equals, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			tc.test(c)
		})
	}
}

func TestGetMux(t *testing.T) {
	c := qt.New(t)

	testCases := []struct {
		name           string
		endpoint       string
		method         string
		hasUser        bool
		isSuperuser    bool
		expectedStatus int
		setupMock      func(service *mock.MockCredentials)
	}{
		{
			name:           "credentials success",
			endpoint:       "/credentials",
			method:         http.MethodGet,
			hasUser:        true,
			expectedStatus: http.StatusOK,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().GetCredentials(gomock.Any()).Return(nil, nil, nil)
			},
		},
		{
			name:           "credentials no user",
			endpoint:       "/credentials",
			method:         http.MethodGet,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials google success",
			endpoint:       "/credentials/google",
			method:         http.MethodPost,
			hasUser:        true,
			isSuperuser:    true,
			expectedStatus: http.StatusOK,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsGoogle(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.CredentialsGoogle{}, nil)
			},
		},
		{
			name:           "credentials google no user",
			endpoint:       "/credentials/google",
			method:         http.MethodPost,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials google not superuser",
			endpoint:       "/credentials/google",
			method:         http.MethodPost,
			hasUser:        true,
			isSuperuser:    false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials deepl success",
			endpoint:       "/credentials/deepl",
			method:         http.MethodPost,
			hasUser:        true,
			isSuperuser:    true,
			expectedStatus: http.StatusOK,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().AddCredentialsDeepL(gomock.Any(), gomock.Any()).Return(&model.CredentialsDeepL{}, nil)
			},
		},
		{
			name:           "credentials deepl no user",
			endpoint:       "/credentials/deepl",
			method:         http.MethodPost,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials deepl not superuser",
			endpoint:       "/credentials/deepl",
			method:         http.MethodPost,
			hasUser:        true,
			isSuperuser:    false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials google remove success",
			endpoint:       "/credentials/google/1",
			method:         http.MethodDelete,
			hasUser:        true,
			isSuperuser:    true,
			expectedStatus: http.StatusNoContent,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsGoogle(gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name:           "credentials google remove no user",
			endpoint:       "/credentials/google/1",
			method:         http.MethodDelete,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials google remove not superuser",
			endpoint:       "/credentials/google/1",
			method:         http.MethodDelete,
			hasUser:        true,
			isSuperuser:    false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials deepl remove success",
			endpoint:       "/credentials/deepl/1",
			method:         http.MethodDelete,
			hasUser:        true,
			isSuperuser:    true,
			expectedStatus: http.StatusNoContent,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().RemoveCredentialsDeepL(gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name:           "credentials deepl remove no user",
			endpoint:       "/credentials/deepl/1",
			method:         http.MethodDelete,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "credentials deepl remove not superuser",
			endpoint:       "/credentials/deepl/1",
			method:         http.MethodDelete,
			hasUser:        true,
			isSuperuser:    false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "sessions google url success",
			endpoint:       "/sessions/google/1?redirect_url=redirectURL",
			method:         http.MethodGet,
			hasUser:        true,
			expectedStatus: http.StatusOK,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionGoogleURL(gomock.Any(), uint(1), gomock.Any(), gomock.Any()).Return("url", nil)
			},
		},
		{
			name:           "sessions google url no user",
			endpoint:       "/sessions/google/1?redirect_url=redirectURL",
			method:         http.MethodGet,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "sessions google",
			endpoint:       "/sessions/google",
			method:         http.MethodGet,
			hasUser:        true,
			expectedStatus: http.StatusOK,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().GetSessionsGoogleByUser(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:           "sessions google no user",
			endpoint:       "/sessions/google",
			method:         http.MethodGet,
			hasUser:        false,
			expectedStatus: http.StatusForbidden,
			setupMock:      func(service *mock.MockCredentials) {},
		},
		{
			name:           "sessions google callback success",
			endpoint:       "/sessions/google/callback",
			method:         http.MethodGet,
			expectedStatus: http.StatusFound,
			setupMock: func(service *mock.MockCredentials) {
				service.EXPECT().CreateSessionGoogle(gomock.Any(), gomock.Any(), gomock.Any()).Return("http://example.com", nil)
			},
		},
	}

	for _, tc := range testCases {
		c.Run(tc.name, func(c *qt.C) {
			ctrl := gomock.NewController(t)
			a := mock.NewMockAuth(ctrl)
			a.EXPECT().AuthMiddleware(gomock.Any()).DoAndReturn(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tc.hasUser {
						r = r.WithContext(auth.ContextWithUser(r.Context(), "userID", tc.isSuperuser))
						next.ServeHTTP(w, r)
					} else {
						http.Error(w, "Forbidden", http.StatusForbidden)
					}
				})
			})

			service := mock.NewMockCredentials(ctrl)
			tc.setupMock(service)

			s := New(service, a, nil, nil)
			mux := s.getMux()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.endpoint, nil)
			mux.ServeHTTP(w, r)

			c.Assert(w.Code, qt.Equals, tc.expectedStatus)
		})
	}
}
