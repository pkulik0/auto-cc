package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/oauth2"
)

const tokenCacheKey = "token_"

func saveToken(rdb *redis.Client, token *oauth2.Token, identity *Identity) error {
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return errors.New("failed to serialize token")
	}

	if err := rdb.Set(tokenCacheKey+identity.Hash(), tokenJson, 0).Err(); err != nil {
		return errors.New(fmt.Sprintf("Failed to save token in db: %s", err.Error()))
	}
	return nil
}

func getTokenFromCache(rdb *redis.Client, identity *Identity) (*oauth2.Token, error) {
	tokenJson, err := rdb.Get(tokenCacheKey + identity.Hash()).Bytes()
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(tokenJson, &token); err != nil {
		return nil, err
	}

	return token, nil
}

func exchangeToken(ctx context.Context, config *oauth2.Config, code string, state string) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

type TokenWrapper struct {
	rdb         *redis.Client
	tokenSource oauth2.TokenSource
	identity    *Identity
}

func NewTokenWrapper(rdb *redis.Client, tokenSource oauth2.TokenSource, identity *Identity) *TokenWrapper {
	return &TokenWrapper{
		rdb:         rdb,
		tokenSource: tokenSource,
		identity:    identity,
	}
}

func (w *TokenWrapper) Token() (*oauth2.Token, error) {
	token, err := w.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if err := saveToken(w.rdb, token, w.identity); err != nil {
		return nil, err
	}

	return token, nil
}
