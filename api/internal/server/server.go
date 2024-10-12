package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/service"
	"github.com/pkulik0/autocc/api/internal/version"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type server struct {
	service service.Service
	auth    auth.Auth
}

func New(s service.Service, a auth.Auth) *server {
	return &server{
		service: s,
		auth:    a,
	}
}

func writeOrLog(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func (s *server) handlerRoot(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(version.Information())
	if err != nil {
		http.Error(w, "failed to read version information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeOrLog(w, data)
}

func (s *server) handlerCredentials(w http.ResponseWriter, r *http.Request) {
	writeOrLog(w, []byte("[]"))
}

func (s *server) Start(port int16) error {
	authMux := http.NewServeMux()
	authMux.HandleFunc("GET /credentials", s.handlerCredentials)

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.handlerRoot)
	mux.Handle("/", s.auth.AuthMiddleware(authMux))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	addr := fmt.Sprintf(":%d", port)
	log.Info().Str("address", addr).Msg("starting server")
	return http.ListenAndServe(addr, logMiddleware(c.Handler(mux)))
}
