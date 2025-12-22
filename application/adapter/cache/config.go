package cache

import (
	"time"

	"src/core/validator"
)

type CacheConfig struct {
	RedisURI  string
	ShortTTL  time.Duration
	MediumTTL time.Duration
	LongTTL   time.Duration
}

var _ validator.IValidable = (*CacheConfig)(nil)

func (c *CacheConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.RedisURI).Required().Default("redis://localhost:6379/0"),
		validator.Number(&c.ShortTTL).Required().Default((time.Minute * 15).Seconds()),
		validator.Number(&c.MediumTTL).Required().Default((time.Hour * 4).Seconds()),
		validator.Number(&c.LongTTL).Required().Default((time.Hour * 24).Seconds()),
	).Validate()
}
