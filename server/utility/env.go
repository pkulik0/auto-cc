package utility

import (
	"github.com/gofiber/fiber/v2/log"
	"os"
)

func GetReqEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set!", key)
	}
	return value
}
