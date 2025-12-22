package cache

import (
	"context"
	"time"

	adapter "src/application/adapter/cache"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/cache/redis"
)

func init() {
	di.RegisterAs[adapter.ICacheAdapter](func() adapter.ICacheAdapter {
		config := &adapter.CacheConfig{
			RedisURI:  env.Get("CACHE_REDIS_URI", "redis://localhost:6379/0"),
			ShortTTL:  time.Minute * 15,
			MediumTTL: time.Hour * 4,
			LongTTL:   time.Hour * 24,
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		impl := impl.NewRedisCacheAdapter(config)
		if err := impl.Ping(context.Background()); err != nil {
			panic(err)
		}
		return impl
	})
}
