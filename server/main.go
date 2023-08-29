package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkulik0/autocc/server/translation"
	"github.com/pkulik0/autocc/server/translation/deepl"
	"github.com/pkulik0/autocc/server/utility"
	"github.com/pkulik0/autocc/server/youtube"
	"os"
)

func setupRedis() *redis.Client {
	url := utility.GetReqEnv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL")
	}

	log.Debugf("Connecting to redis on %s", opts.Addr)
	return redis.NewClient(opts)
}

func main() {
	port := utility.GetReqEnv("PORT")
	apiKey := utility.GetReqEnv("DEEPL_API_KEY")
	apiUrl := os.Getenv("DEEPL_API_URL")
	if apiUrl == "" {
		apiUrl = "https://api-free.deepl.com/v2"
		log.Infof("DEEPL_API_URL not set. Using default: %s", apiUrl)
	}

	rdb := setupRedis()
	defer rdb.Close()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	ytService := youtube.NewYoutubeService(rdb)
	ytService.RegisterEndpoints(app)

	deeplClient := deepl.NewDeepL(apiUrl, apiKey)
	translationService := translation.NewTranslationService(deeplClient, rdb)
	translationService.RegisterEndpoints(app)

	log.Infof("Listening on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to serve app: %s", err)
	}
}
