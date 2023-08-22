package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
)

type Service struct {
	rdb          *redis.Client
	oauth2Config *oauth2.Config
}

func newService(rdb *redis.Client, oauth2Config *oauth2.Config) Service {
	return Service{
		rdb:          rdb,
		oauth2Config: oauth2Config,
	}
}

func (s *Service) registerEndpoints(app *fiber.App) {
	app.Get("/auth", s.authHandler)
	app.Get("/callback", s.callbackHandler)
	app.Use(s.authMiddleware)

	app.Get("/test", func(ctx *fiber.Ctx) error {
		return ctx.SendString("TEST")
	})
}

func (s *Service) callbackHandler(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	if code == "" || state == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing oauth2 data.")
	}

	token, err := exchangeToken(s.oauth2Config, code, state)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to exchange token.")
	}

	if err := saveToken(s.rdb, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save token.")
	}

	return ctx.SendString("OK")
}

func (s *Service) authHandler(ctx *fiber.Ctx) error {
	url := s.oauth2Config.AuthCodeURL("stateTODO", oauth2.AccessTypeOffline)
	return ctx.Redirect(url, fiber.StatusFound)
}

func (s *Service) authMiddleware(ctx *fiber.Ctx) error {
	token, err := getTokenFromCache(s.rdb)
	if err != nil {
		log.Error(err)
		return ctx.Status(fiber.StatusUnauthorized).SendString("No oauth2 token found.")
	}

	log.Info("?")
	ctx.Locals("token", token)
	return ctx.Next()
}
