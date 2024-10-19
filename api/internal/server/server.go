package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/service"
	"github.com/pkulik0/autocc/api/internal/version"
	"github.com/pkulik0/autocc/api/internal/youtube"
)

type GetYoutubeServiceFunc func(string) youtube.Youtube

type server struct {
	service service.Service
	auth    auth.Auth
	youtube youtube.Youtube
}

func New(service service.Service, auth auth.Auth, youtube youtube.Youtube) *server {
	return &server{
		service: service,
		auth:    auth,
		youtube: youtube,
	}
}

func (s *server) handlerRoot(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(version.Information())
	if err != nil {
		errLog(w, err, "failed to marshal version information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeOrLog(w, data)
}

func (s *server) handlerCredentials(w http.ResponseWriter, r *http.Request) {
	credentialsGoogle, credentialsDeepL, err := s.service.GetCredentials(r.Context())
	if err != nil {
		errLog(w, err, "failed to get credentials", http.StatusInternalServerError)
		return
	}

	var resp pb.GetCredentialsResponse
	for _, c := range credentialsGoogle {
		resp.Google = append(resp.Google, c.ToProto())
	}
	for _, c := range credentialsDeepL {
		resp.Deepl = append(resp.Deepl, c.ToProto())
	}

	writePb(w, &resp)
}

func (s *server) handlerAddCredentialsGoogle(w http.ResponseWriter, r *http.Request) {
	var req pb.AddCredentialsGoogleRequest
	err := readPb(r, &req)
	if err != nil {
		errLog(w, err, "failed to decode request", http.StatusBadRequest)
		return
	}

	credentials, err := s.service.AddCredentialsGoogle(r.Context(), req.ClientId, req.ClientSecret)
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to create google credentials", http.StatusInternalServerError)
		return
	}

	writePb(w, &pb.AddCredentialsGoogleResponse{Credentials: credentials.ToProto()})
}

func (s *server) handlerAddCredentialsDeepL(w http.ResponseWriter, r *http.Request) {
	var req pb.AddCredentialsDeepLRequest
	err := readPb(r, &req)
	if err != nil {
		errLog(w, err, "failed to decode request", http.StatusBadRequest)
		return
	}

	credentials, err := s.service.AddCredentialsDeepL(r.Context(), req.Key)
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to create deepl credentials", http.StatusInternalServerError)
		return
	}

	writePb(w, &pb.AddCredentialsDeepLResponse{Credentials: credentials.ToProto()})
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
		errLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	err = s.service.RemoveCredentialsGoogle(r.Context(), uint(id))
	if err != nil {
		errLog(w, err, "failed to remove google credentials", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerRemoveCredentialsDeepL(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathID(r)
	if err != nil {
		errLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	err = s.service.RemoveCredentialsDeepL(r.Context(), uint(id))
	if err != nil {
		errLog(w, err, "failed to remove deepl credentials", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerSessionGoogleURL(w http.ResponseWriter, r *http.Request) {
	credentialsID, err := parsePathID(r)
	if err != nil {
		errLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	redirectURL := r.URL.Query().Get("redirect_url")

	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	url, err := s.service.GetSessionGoogleURL(r.Context(), credentialsID, userID, redirectURL)
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to get google session url", http.StatusInternalServerError)
		return
	}

	var req pb.GetSessionGoogleURLResponse
	req.Url = url
	writePb(w, &req)
}

func (s *server) handlerSessionGoogleCallback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errLog(w, err, "failed to parse form", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	code := r.FormValue("code")

	url, err := s.service.CreateSessionGoogle(r.Context(), state, code)
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to create google session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *server) handlerUserSessionsGoogle(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	sessions, err := s.service.GetSessionsGoogleByUser(r.Context(), userID)
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to get google sessions", http.StatusInternalServerError)
		return
	}

	var resp pb.GetUserSessionsGoogleResponse
	for _, s := range sessions {
		resp.CredentialIds = append(resp.CredentialIds, uint64(s.CredentialsID))
	}
	writePb(w, &resp)
}

func (s *server) handlerRemoveSessionGoogle(w http.ResponseWriter, r *http.Request) {
	credentialsID, err := parsePathID(r)
	if err != nil {
		errLog(w, err, "failed to parse id", http.StatusBadRequest)
		return
	}

	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	err = s.service.RemoveSessionGoogle(r.Context(), userID, uint(credentialsID))
	switch err {
	case nil:
	case service.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to remove google session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerYoutubeVideos(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	nextPageToken := r.URL.Query().Get("next_page_token")
	videos, nextPageToken, err := s.youtube.GetVideos(r.Context(), userID, nextPageToken)
	switch err {
	case nil:
	case youtube.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to get videos", http.StatusInternalServerError)
		return
	}

	var resp pb.GetYoutubeVideosResponse
	resp.Videos = videos
	resp.NextPageToken = nextPageToken

	writePb(w, &resp)
}

func (s *server) handlerYoutubeMetadata(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	metadata, err := s.youtube.GetMetadata(r.Context(), userID, id)
	switch err {
	case nil:
	case youtube.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	case youtube.ErrNotFound:
		errLog(w, err, "video not found", http.StatusNotFound)
		return
	default:
		errLog(w, err, "failed to get video metadata", http.StatusInternalServerError)
	}

	var resp pb.GetMetadataResponse
	resp.Metadata = metadata

	writePb(w, &resp)
}

func (s *server) handlerYoutubeUpdateMetadata(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	var req pb.UpdateMetadataRequest
	err := readPb(r, &req)
	if err != nil {
		errLog(w, err, "failed to decode request", http.StatusBadRequest)
		return
	}

	err = s.youtube.UpdateMetadata(r.Context(), userID, id, req.Metadata)
	switch err {
	case nil:
	case youtube.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to update video metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handlerYoutubeCC(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	cc, err := s.youtube.GetClosedCaptions(r.Context(), userID, id)
	switch err {
	case nil:
	case youtube.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to get video cc", http.StatusInternalServerError)
	}

	var resp pb.GetClosedCaptionsResponse
	resp.ClosedCaptions = cc

	writePb(w, &resp)
}

func (s *server) handlerYoutubeDownloadCC(w http.ResponseWriter, r *http.Request) {
	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	srt, err := s.youtube.DownloadClosedCaptions(r.Context(), userID, id)
	switch err {
	case nil:
	case youtube.ErrInvalidInput:
		errLog(w, err, "invalid input", http.StatusBadRequest)
		return
	default:
		errLog(w, err, "failed to download video cc", http.StatusInternalServerError)
	}

	var resp pb.DownloadClosedCaptionsResponse
	resp.Srt = srt

	writePb(w, &resp)
}

func superuserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, isSuperuser, ok := auth.UserFromContext(r.Context())
		if !ok {
			errLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
			return
		}
		if !isSuperuser {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) getMux() *http.ServeMux {
	superuserMux := http.NewServeMux()
	superuserMux.HandleFunc("POST /credentials/google", s.handlerAddCredentialsGoogle)
	superuserMux.HandleFunc("POST /credentials/deepl", s.handlerAddCredentialsDeepL)
	superuserMux.HandleFunc("DELETE /credentials/google/{id}", s.handlerRemoveCredentialsGoogle)
	superuserMux.HandleFunc("DELETE /credentials/deepl/{id}", s.handlerRemoveCredentialsDeepL)

	ytMux := http.NewServeMux()
	ytMux.HandleFunc("GET /videos", s.handlerYoutubeVideos)
	ytMux.HandleFunc("GET /videos/{id}/metadata", s.handlerYoutubeMetadata)
	ytMux.HandleFunc("PUT /videos/{id}/metadata", s.handlerYoutubeUpdateMetadata)
	ytMux.HandleFunc("GET /videos/{id}/cc", s.handlerYoutubeCC)
	ytMux.HandleFunc("GET /cc/{id}", s.handlerYoutubeDownloadCC)

	authMux := http.NewServeMux()
	authMux.Handle("/", superuserMiddleware(superuserMux))
	authMux.HandleFunc("GET /credentials", s.handlerCredentials)
	authMux.HandleFunc("GET /sessions/google", s.handlerUserSessionsGoogle)
	authMux.HandleFunc("GET /sessions/google/{id}", s.handlerSessionGoogleURL)
	authMux.HandleFunc("DELETE /sessions/google/{id}", s.handlerRemoveSessionGoogle)
	authMux.Handle("/youtube/", http.StripPrefix("/youtube", ytMux))

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.handlerRoot)
	mux.Handle("/", s.auth.AuthMiddleware(authMux))
	mux.HandleFunc("GET /sessions/google/callback", s.handlerSessionGoogleCallback)

	return mux
}

func (s *server) Start(port uint16) error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	addr := fmt.Sprintf(":%d", port)
	log.Info().Str("address", addr).Msg("starting server")
	return http.ListenAndServe(addr, logMiddleware(c.Handler(s.getMux())))
}
