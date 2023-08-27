package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"time"
)

type Service struct {
	rdb      *redis.Client
	identity *Identity
}

func (s *Service) Identity() (*Identity, error) {
	if s.identity == nil {
		identity, err := s.nextIdentity()
		if err != nil {
			return nil, err
		}
		return identity, nil
	}
	return s.identity, nil
}

func newService(rdb *redis.Client) Service {
	return Service{
		rdb: rdb,
	}
}

func (s *Service) registerEndpoints(app *fiber.App) {
	app.Get("/identities", s.getIdentitiesHandler)
	app.Post("/identities", s.addIdentityHandler)

	app.Get("/auth", s.authUrlsHandler)
	app.Get("/callback", s.callbackHandler)

	app.Use(s.authMiddleware)
	app.Use(cache.New(cache.Config{
		Expiration: time.Minute * 15,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.Path() + ctx.Query("token")
		},
	}))

	app.Get("/videos", s.videosHandler)
	app.Get("/videos/:videoId/cc", s.ccListHandler)
	app.Post("/videos/:videoId/cc", s.ccUpload)
	app.Get("/cc/:ccId", s.ccDownloadHandler)
}
