package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if c == nil || c.client == nil {
		return nil, nil
	}

	value, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Close()
}
