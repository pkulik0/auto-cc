package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (s *Service) authMiddleware(ctx *fiber.Ctx) error {
	token, err := getTokenFromCache(s.rdb)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("No access.")
	}

	tokenSource := s.oauth2Config.TokenSource(ctx.Context(), token)
	tokenWrapper := NewTokenWrapper(s.rdb, tokenSource)

	oauth2Client := oauth2.NewClient(ctx.Context(), tokenWrapper)
	youtubeClient, err := youtube.NewService(ctx.Context(), option.WithHTTPClient(oauth2Client))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create youtube client.")
	}

	ctx.Locals("youtubeClient", youtubeClient)
	return ctx.Next()
}
