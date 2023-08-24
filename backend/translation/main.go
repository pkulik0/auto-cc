package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkulik0/autocc/translation/deepl"
	"os"
)

func getReqEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set!", key)
	}
	return value
}

func setupRedis() *redis.Client {
	url := getReqEnv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL")
	}

	log.Debugf("Connecting to redis on %s", opts.Addr)
	return redis.NewClient(opts)
}

func main() {
	port := getReqEnv("PORT")
	apiKey := getReqEnv("DEEPL_API_KEY")
	apiUrl := os.Getenv("DEEPL_API_URL")
	if apiUrl == "" {
		apiUrl = "https://api-free.deepl.com/v2"
		log.Infof("DEEPL_API_URL not set. Using default: %s", apiUrl)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	deeplClient := deepl.NewDeepL(apiUrl, apiKey)
	rdb := setupRedis()

	service := newService(deeplClient, rdb)
	service.registerEndpoint(app)

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}
}
