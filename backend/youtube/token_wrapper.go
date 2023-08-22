package main

import (
	"github.com/go-redis/redis"
	"golang.org/x/oauth2"
)

type TokenWrapper struct {
	rdb         *redis.Client
	tokenSource oauth2.TokenSource
}

func NewTokenWrapper(rdb *redis.Client, tokenSource oauth2.TokenSource) *TokenWrapper {
	return &TokenWrapper{
		rdb:         rdb,
		tokenSource: tokenSource,
	}
}

func (w *TokenWrapper) Token() (*oauth2.Token, error) {
	token, err := w.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if err := saveToken(w.rdb, token); err != nil {
		return nil, err
	}

	return token, nil
}
