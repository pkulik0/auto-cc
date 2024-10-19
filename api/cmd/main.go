package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/credentials"
	"github.com/pkulik0/autocc/api/internal/oauth"
	"github.com/pkulik0/autocc/api/internal/server"
	"github.com/pkulik0/autocc/api/internal/store"
	"github.com/pkulik0/autocc/api/internal/translation"
	"github.com/pkulik0/autocc/api/internal/version"
	"github.com/pkulik0/autocc/api/internal/youtube"
)

func main() {
	version.EnsureSet()

	c, err := parseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Second
	if c.IsProduction {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("running in production mode")
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Info().Msg("running in development mode")
	}

	log.Info().Str("version", version.Version).Str("build_time", version.BuildTime).Msg("AutoCC API")

	auth, err := auth.New(context.Background(), c.KeycloakURL, c.KeycloakRealm, c.KeycloakClientId, c.KeycloakClientSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create auth")
	}

	store, err := store.New(c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPass, c.PostgresDB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create store")
	}

	translator := translation.New(store)
	credentials := credentials.New(store, oauth.New(c.GoogleCallbackURL), translator)

	server := server.New(credentials, auth, youtube.New(store), translator)
	err = server.Start(c.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
