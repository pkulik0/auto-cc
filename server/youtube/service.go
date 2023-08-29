package youtube

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

func NewYoutubeService(rdb *redis.Client) Service {
	return Service{
		rdb: rdb,
	}
}

func (s *Service) RegisterEndpoints(app *fiber.App) {
	group := app.Group("/youtube")

	group.Get("/identities", s.getIdentitiesHandler)
	group.Post("/identities", s.addIdentityHandler)

	group.Get("/auth", s.authUrlsHandler)
	group.Get("/callback", s.callbackHandler)

	group.Use(s.authMiddleware)
	group.Use(cache.New(cache.Config{
		Expiration: time.Minute * 15,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.Path() + ctx.Query("token")
		},
	}))

	group.Get("/videos", s.videosHandler)

	group.Get("/videos/:videoId", s.videoMetadataHandler)
	group.Post("/videos/:videoId", s.videoUpdateMetadataHandler)

	group.Get("/videos/:videoId/cc", s.ccListHandler)
	group.Post("/videos/:videoId/cc", s.ccUpload)

	group.Get("/cc/:ccId", s.ccDownloadHandler)
}
