package youtube

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Service) authMiddleware(ctx *fiber.Ctx) error {
	if err := s.generateClientFromCurrentIdentity(ctx); err != nil {
		return err
	}
	return ctx.Next()
}
