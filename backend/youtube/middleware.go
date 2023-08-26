package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (s *Service) authMiddleware(ctx *fiber.Ctx) error {
	identity, err := s.Identity()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Identity error.")
	}

	token, err := getTokenFromCache(s.rdb, identity)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("No access.")
	}

	config := identity.getOAuth2Config()
	tokenSource := config.TokenSource(ctx.Context(), token)
	tokenWrapper := NewTokenWrapper(s.rdb, tokenSource, identity)

	oauth2Client := oauth2.NewClient(ctx.Context(), tokenWrapper)
	youtubeClient, err := youtube.NewService(ctx.Context(), option.WithHTTPClient(oauth2Client))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create youtube client.")
	}

	ctx.Locals("youtubeClient", youtubeClient)
	return ctx.Next()
}
