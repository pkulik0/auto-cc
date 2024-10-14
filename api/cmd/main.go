package main

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/server"
	"github.com/pkulik0/autocc/api/internal/service"
	"github.com/pkulik0/autocc/api/internal/store"
	"github.com/pkulik0/autocc/api/internal/version"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Second
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	version.EnsureSet()
	log.Info().Str("version", version.Version).Str("build_time", version.BuildTime).Msg("AutoCC API")

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	c, err := parseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}

	store, err := store.New(c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPass, c.PostgresDB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create store")
	}
	service := service.New(store, c.GoogleCallbackURL)

	auth, err := auth.New(context.Background(), c.KeycloakURL, c.KeycloakRealm, c.KeycloakClientId, c.KeycloakClientSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create auth")
	}

	server := server.New(service, auth, c.GoogleRedirectURL)
	err = server.Start(c.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
