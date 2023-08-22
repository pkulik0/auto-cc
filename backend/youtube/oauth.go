package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const tokenCacheKey string = "oauth2-token"

func getOAuth2Config() *oauth2.Config {
	clientId := getEnv("GOOGLE_CLIENT_ID")
	clientSecret := getEnv("GOOGLE_CLIENT_SECRET")
	redirectUri := getEnv("GOOGLE_REDIRECT_URI")

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUri,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint:     google.Endpoint,
	}
}

func saveToken(rdb *redis.Client, token *oauth2.Token) error {
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return errors.New("failed to serialize token")
	}

	if err := rdb.Set(tokenCacheKey, tokenJson, 0).Err(); err != nil {
		return errors.New(fmt.Sprintf("Failed to save token in db: %s", err.Error()))
	}
	return nil
}

func getTokenFromCache(rdb *redis.Client) (*oauth2.Token, error) {
	tokenJson, err := rdb.Get(tokenCacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(tokenJson, &token); err != nil {
		return nil, err
	}

	return token, nil
}

func exchangeToken(config *oauth2.Config, code string, state string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func getNewToken(config *oauth2.Config) *oauth2.Token {
	url := config.AuthCodeURL("stateTODO", oauth2.AccessTypeOffline)
	log.Infof("URL: %s", url)

	return nil
}
