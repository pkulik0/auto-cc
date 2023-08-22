package main

import (
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkulik0/autocc/translation/deepl"
	"os"
)

func main() {
	apiUrl := os.Getenv("DEEPL_API_URL")
	if apiUrl == "" {
		apiUrl = "https://api-free.deepl.com/v2"
		log.Infof("DEEPL_API_URL not set. Using default: %s", apiUrl)
	}

	apiKey := os.Getenv("DEEPL_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPL_API_KEY not set")
	}

	deeplClient := deepl.NewDeepL(apiUrl, apiKey)

	languages, err := deeplClient.GetLanguages(deepl.TargetLanguages)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(languages)

	usageInfo, err := deeplClient.GetUsage()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(usageInfo)
}
