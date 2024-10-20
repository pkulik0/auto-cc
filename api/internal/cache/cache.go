package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Cache is the interface that wraps cache methods.
//
//go:generate mockgen -destination=../mock/cache.go -package=mock . Cache
type Cache interface {
	// Get returns the value for the key.
	Get(ctx context.Context, key string) (string, error)

	// Set sets the value for the key with expiration.
	Set(ctx context.Context, key, value string, expiration time.Duration) error

	// Del deletes the value for the key.
	Del(ctx context.Context, key string) error
}

var _ Cache = &redisCache{}

type redisCache struct {
	client *redis.Client
}

// New creates a new cache service.
func New(ctx context.Context, addr string) (*redisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	log.Debug().Msg("created new cache instance")
	return &redisCache{
		client: client,
	}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *redisCache) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisCache) Del(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
