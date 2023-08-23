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

func newService(deeplClient *deepl.DeepL, rdb *redis.Client) *Service {
	return &Service{
		deeplClient: deeplClient,
		rdb:         rdb,
	}
}

func (s *Service) registerEndpoint(app *fiber.App) {
	app.Get("/languages", s.languagesHandler)
	app.Get("/translate", s.translateHandler)
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

type translateRequest struct {
	SourceLanguageCode string   `json:"source"`
	TargetLanguageCode string   `json:"target"`
	Text               []string `json:"text"`
}

func (r *translateRequest) isValid() bool {
	return len(r.Text) > 0 && len(r.TargetLanguageCode) >= 2 && len(r.SourceLanguageCode) >= 2
}

func (s *Service) translateHandler(ctx *fiber.Ctx) error {
	var requestData translateRequest
	err := ctx.BodyParser(&requestData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body.")
	}
	if !requestData.isValid() {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request data.")
	}

	translatedText, err := s.deeplClient.Translate(requestData.Text, requestData.SourceLanguageCode, requestData.TargetLanguageCode)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch translations.")
	}

	return ctx.JSON(translatedText)
}
