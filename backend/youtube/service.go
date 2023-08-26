package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
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

	app.Get("/videos", s.videosHandler)
	app.Get("/videos/:videoId/cc", s.ccListHandler)
	app.Post("/videos/:videoId/cc", s.ccUpload)
	app.Get("/cc/:ccId", s.ccDownloadHandler)
}
