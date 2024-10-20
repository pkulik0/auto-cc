package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/autocc"
	"github.com/pkulik0/autocc/api/internal/cache"
	"github.com/pkulik0/autocc/api/internal/credentials"
	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/helpers"
	"github.com/pkulik0/autocc/api/internal/middleware"
	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/version"
	"github.com/pkulik0/autocc/api/internal/youtube"
)

type server struct {
	cache       cache.Cache
	credentials credentials.Credentials
	auth        auth.Auth
	youtube     youtube.Youtube
	autocc      autocc.AutoCC
}

func New(cache cache.Cache, credentials credentials.Credentials, auth auth.Auth, youtube youtube.Youtube, autocc autocc.AutoCC) *server {
	return &server{
		cache:       cache,
		credentials: credentials,
		auth:        auth,
		youtube:     youtube,
		autocc:      autocc,
	}
}

func (s *server) handlerRoot(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(version.Information())
	if err != nil {
		helpers.ErrLog(w, err, "failed to marshal version information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteOrLog(w, data)
}

func (s *server) handlerCredentials(w http.ResponseWriter, r *http.Request) {
	credentialsGoogle, credentialsDeepL, err := s.credentials.GetCredentials(r.Context())
	if err != nil {
		helpers.ErrLog(w, err, "failed to get credentials", http.StatusInternalServerError)
		return
	}

	var resp pb.GetCredentialsResponse
	for _, c := range credentialsGoogle {
		resp.Google = append(resp.Google, c.ToProto())
	}
	for _, c := range credentialsDeepL {
		resp.Deepl = append(resp.Deepl, c.ToProto())
	}

	helpers.WritePb(w, &resp)
}

func (s *server) handlerAddCredentialsGoogle(w http.ResponseWriter, r *http.Request) {
	var req pb.AddCredentialsGoogleRequest
	err := helpers.ReadPb(r, &req)
	if err != nil {
		helpers.ErrLog(w, err, "failed to decode request", http.StatusBadRequest)
		return
	}

	cred, err := s.credentials.AddCredentialsGoogle(r.Context(), req.ClientId, req.ClientSecret)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to create google credentials", http.StatusInternalServerError)
		return
	}

	helpers.WritePb(w, &pb.AddCredentialsGoogleResponse{Credentials: cred.ToProto()})
}

func (s *server) handlerAddCredentialsDeepL(w http.ResponseWriter, r *http.Request) {
	var req pb.AddCredentialsDeepLRequest
	err := helpers.ReadPb(r, &req)
	if err != nil {
		helpers.ErrLog(w, err, "failed to decode request", http.StatusBadRequest)
		return
	}

	cred, err := s.credentials.AddCredentialsDeepL(r.Context(), req.Key)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to create deepl credentials", http.StatusInternalServerError)
		return
	}

	helpers.WritePb(w, &pb.AddCredentialsDeepLResponse{Credentials: cred.ToProto()})
}

func parsePathID(r *http.Request) (uint, error) {
	idValue := r.PathValue("id")
	if idValue == "" {
		return 0, fmt.Errorf("missing id")
	}

	id, err := strconv.ParseUint(idValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse id")
	}

	return uint(id), nil
}

func (s *server) handlerRemoveCredentialsGoogle(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r)
	if err != nil {
		helpers.ErrLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	err = s.credentials.RemoveCredentialsGoogle(r.Context(), uint(id))
	if err != nil {
		helpers.ErrLog(w, err, "failed to remove google credentials", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerRemoveCredentialsDeepL(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r)
	if err != nil {
		helpers.ErrLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	err = s.credentials.RemoveCredentialsDeepL(r.Context(), uint(id))
	if err != nil {
		helpers.ErrLog(w, err, "failed to remove deepl credentials", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerSessionGoogleURL(w http.ResponseWriter, r *http.Request) {
	credentialsID, err := parsePathID(r)
	if err != nil {
		helpers.ErrLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	redirectURL := r.URL.Query().Get("redirect_url")

	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	url, err := s.credentials.GetSessionGoogleURL(r.Context(), credentialsID, userID, redirectURL)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to get google session url", http.StatusInternalServerError)
		return
	}

	var req pb.GetSessionGoogleURLResponse
	req.Url = url
	helpers.WritePb(w, &req)
}

func (s *server) handlerSessionGoogleCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ErrLog(w, err, "failed to parse form", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	code := r.FormValue("code")

	url, err := s.credentials.CreateSessionGoogle(r.Context(), state, code)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to create google session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *server) handlerUserSessionsGoogle(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	sessions, err := s.credentials.GetSessionsGoogleByUser(r.Context(), userID)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to get google sessions", http.StatusInternalServerError)
		return
	}

	var resp pb.GetUserSessionsGoogleResponse
	for _, s := range sessions {
		resp.CredentialIds = append(resp.CredentialIds, uint64(s.CredentialsID))
	}
	helpers.WritePb(w, &resp)
}

func (s *server) handlerRemoveSessionGoogle(w http.ResponseWriter, r *http.Request) {
	credentialsID, err := parsePathID(r)
	if err != nil {
		helpers.ErrLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	err = s.credentials.RemoveSessionGoogle(r.Context(), userID, uint(credentialsID))
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to remove google session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerYoutubeVideos(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	nextPageToken := r.URL.Query().Get("next_page_token")
	videos, nextPageToken, err := s.youtube.GetVideos(r.Context(), userID, nextPageToken)
	switch err {
	case nil:
	case errs.NotFound:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		helpers.ErrLog(w, err, "failed to get videos", http.StatusInternalServerError)
		return
	}

	var resp pb.GetYoutubeVideosResponse
	resp.Videos = videos
	resp.NextPageToken = nextPageToken

	helpers.WritePb(w, &resp)
}

func (s *server) handlerProcess(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	videoID := r.PathValue("id")

	err := s.autocc.Process(r.Context(), userID, videoID)
	switch err {
	case nil:
	case errs.InvalidInput:
		helpers.ErrLog(w, err, "invalid input", http.StatusBadRequest)
		return
	case errs.SourceClosedCaptionsNotFound:
		helpers.ErrLog(w, err, "source closed captions not found", http.StatusNotFound)
		return
	case errs.NotFound:
		helpers.ErrLog(w, err, "not found", http.StatusNotFound)
		return
	default:
		helpers.ErrLog(w, err, "failed to process video", http.StatusInternalServerError)
		return
	}

	log.Info().Str("video_id", videoID).Msg("finished processing video")
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getMux() *http.ServeMux {
	superuserMux := http.NewServeMux()
	superuserMux.HandleFunc("POST /credentials/google", s.handlerAddCredentialsGoogle)
	superuserMux.HandleFunc("POST /credentials/deepl", s.handlerAddCredentialsDeepL)
	superuserMux.HandleFunc("DELETE /credentials/google/{id}", s.handlerRemoveCredentialsGoogle)
	superuserMux.HandleFunc("DELETE /credentials/deepl/{id}", s.handlerRemoveCredentialsDeepL)

	ytMux := http.NewServeMux()
	ytMux.HandleFunc("GET /videos", s.handlerYoutubeVideos)
	ytMux.HandleFunc("POST /videos/{id}", s.handlerProcess)

	authMux := http.NewServeMux()
	authMux.Handle("/", middleware.Superuser(superuserMux))
	authMux.HandleFunc("GET /credentials", s.handlerCredentials)
	authMux.HandleFunc("GET /sessions/google", s.handlerUserSessionsGoogle)
	authMux.HandleFunc("GET /sessions/google/{id}", s.handlerSessionGoogleURL)
	authMux.HandleFunc("DELETE /sessions/google/{id}", s.handlerRemoveSessionGoogle)
	authMux.Handle("/youtube/", http.StripPrefix("/youtube", ytMux))

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.handlerRoot)
	mux.Handle("/", middleware.Auth(s.auth, authMux))
	mux.HandleFunc("GET /sessions/google/callback", s.handlerSessionGoogleCallback)

	return mux
}

// Start starts the server on the given port.
func (s *server) Start(port uint16) error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	addr := fmt.Sprintf(":%d", port)
	log.Info().Str("address", addr).Msg("starting server")
	return http.ListenAndServe(addr, http.TimeoutHandler(middleware.Panic(middleware.Log(c.Handler(s.getMux()))), 15*time.Second, ""))
}
