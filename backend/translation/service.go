package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/pkulik0/autocc/translation/deepl"
)

type Service struct {
	deeplClient *deepl.DeepL
	rdb         *redis.Client
}

func NewService(deeplClient *deepl.DeepL, rdb *redis.Client) *Service {
	return &Service{
		deeplClient: deeplClient,
		rdb:         rdb,
	}
}

func (s *Service) registerEndpoint(app *fiber.App) {
	app.Get("/languages", s.languagesHandler)
}

func (s *Service) languagesHandler(ctx *fiber.Ctx) error {
	targetLangs, err := s.deeplClient.GetLanguages(deepl.TargetLanguages)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch target languages.")
	}

	sourceLangs, err := s.deeplClient.GetLanguages(deepl.SourceLanguages)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch source languages.")
	}

	return ctx.JSON(fiber.Map{
		"target": targetLangs,
		"source": sourceLangs,
	})
}
