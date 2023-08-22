package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set!", key)
	}
	return value
}

func setupRedis() *redis.Client {
	url := getEnv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL")
	}

	log.Debugf("Connecting to redis on %s", opts.Addr)
	return redis.NewClient(opts)
}

func main() {
	log.Info("Starting...")
	service := newService(setupRedis(), getOAuth2Config())

	app := fiber.New()
	app.Use(logger.New())
	service.registerEndpoints(app)

	port := getEnv("PORT")
	log.Infof("Listening on port %s", port)
	addr := fmt.Sprintf(":%s", port)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to serve app: %s", err.Error())
	}

	log.Info("Bye")
}
