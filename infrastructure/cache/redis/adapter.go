package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"

	adapter "src/application/adapter/cache"
)

type RedisCacheAdapter struct {
	client *goredis.Client
	config *adapter.CacheConfig
}

var _ adapter.ICacheAdapter = (*RedisCacheAdapter)(nil)

func NewRedisCacheAdapter(config *adapter.CacheConfig) *RedisCacheAdapter {
	opt, err := goredis.ParseURL(config.RedisURI)
	if err != nil {
		panic(err)
	}
	client := goredis.NewClient(opt)
	return &RedisCacheAdapter{client: client, config: config}
}

func (a *RedisCacheAdapter) Config() *adapter.CacheConfig {
	return a.config
}

func (a *RedisCacheAdapter) Ping(ctx context.Context) error {
	return a.client.Ping(ctx).Err()
}

func (a *RedisCacheAdapter) Has(ctx context.Context, key string) (bool, error) {
	n, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (a *RedisCacheAdapter) Set(ctx context.Context, key string, value string, optionalTimeToLive ...time.Duration) error {
	var ttl time.Duration
	if len(optionalTimeToLive) > 0 {
		ttl = optionalTimeToLive[0]
	}
	return a.client.Set(ctx, key, value, ttl).Err()
}

func (a *RedisCacheAdapter) Get(ctx context.Context, key string) (string, bool, error) {
	res, err := a.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return res, true, nil
}

func (a *RedisCacheAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}
