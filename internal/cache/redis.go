package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb           *redis.Client
	ctx           context.Context
	defaultExpiry time.Duration
}

func NewRedisCache(rdb *redis.Client, ctx context.Context, defaultExpiry time.Duration) *RedisCache {
	return &RedisCache{rdb: rdb, ctx: ctx, defaultExpiry: defaultExpiry}
}

func ConnectRedis(redisURL string) (*redis.Client, context.Context) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		DB:       0,
		Password: "",
		Protocol: 2,
	})

	return rdb, ctx
}

func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.rdb.Get(c.ctx, key).Result()

	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisCache) Set(key string, val string) error {
	_, err := c.rdb.Set(c.ctx, key, val, c.defaultExpiry).Result()

	return err
}

func (c *RedisCache) Del(key string) error {
	_, err := c.rdb.Del(c.ctx, key).Result()

	return err
}
