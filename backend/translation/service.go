package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkulik0/autocc/translation/deepl"
	"time"
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
	app.Post("/translate", s.translateHandler)
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
	if err := ctx.BodyParser(&requestData); err != nil {
		log.Errorf("Failed to parse request body: %s", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body.")
	}
	if !requestData.isValid() {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request data.")
	}

	cacheKey := fmt.Sprintf("translation_%s_%s_%s", requestData.SourceLanguageCode, requestData.TargetLanguageCode, hashFromStrings(requestData.Text))
	cachedTranslation, err := s.rdb.LRange(cacheKey, 0, -1).Result()
	if err != nil {
		log.Errorf("Failed to get %s from db: %s", cacheKey, err)
	} else if len(cachedTranslation) > 0 {
		return ctx.JSON(cachedTranslation)
	}

	translatedText, err := s.deeplClient.Translate(requestData.Text, requestData.SourceLanguageCode, requestData.TargetLanguageCode)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch translations.")
	}

	if err := s.rdb.RPush(cacheKey, translatedText).Err(); err != nil {
		log.Errorf("Failed to cache translations: %s", err)
	}
	if err := s.rdb.Expire(cacheKey, time.Hour*24).Err(); err != nil {
		log.Errorf("Failed set set expiry time on %s: %s", cacheKey, err)
	}
	return ctx.JSON(translatedText)
}
