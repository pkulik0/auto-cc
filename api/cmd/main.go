package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/server"
	"github.com/pkulik0/autocc/api/internal/service"
	"github.com/pkulik0/autocc/api/internal/version"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	version.EnsureSet()
	log.Info().Msg("AutoCC API started")

	c, err := parseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}

	service := service.New()
	server := server.New(service)

	err = server.Start(c.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
