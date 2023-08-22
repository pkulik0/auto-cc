package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
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

	app.Get("/videos", s.videosHandler)
	app.Get("/videos/:videoId/cc", s.ccListHandler)
}

func (s *Service) callbackHandler(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	if code == "" || state == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing oauth2 code/state.")
	}

	token, err := exchangeToken(ctx.Context(), s.oauth2Config, code, state)
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

func (s *Service) videosHandler(ctx *fiber.Ctx) error {
	youtubeClient, ok := ctx.Locals("youtubeClient").(*youtube.Service)
	if !ok {
		log.Error("Failed to get oauth2Client from locals.")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal error.")
	}

	res, err := youtubeClient.Search.List([]string{"snippet"}).ForMine(true).MaxResults(50).Type("video").Do()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}
	return ctx.JSON(res)
}

func (s *Service) ccListHandler(ctx *fiber.Ctx) error {
	youtubeClient, ok := ctx.Locals("youtubeClient").(*youtube.Service)
	if !ok {
		log.Error("Failed to get oauth2Client from locals.")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal error.")
	}

	videoId := ctx.Params("videoId")
	if videoId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing video id.")
	}

	res, err := youtubeClient.Captions.List([]string{"id"}, videoId).Do()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}

	return ctx.JSON(res)
}
