package cache

import (
	"context"
	"time"
)

// ICacheAdapter is a generic key-value contract with optional TTL.
// Concrete implementations are in infra/cache/*.
type ICacheAdapter interface {
	Config() *CacheConfig

	// Ping is used to check if the connection to the cache is working.
	Ping(ctx context.Context) error

	// Has checks if a key exists in the cache.
	Has(ctx context.Context, key string) (bool, error)

	// Set defines a value for a key.
	// If TimeToLive <= 0, the implementation decides if it expires or not (typically "no expiration").
	Set(ctx context.Context, key string, value string, optionalTimeToLive ...time.Duration) error

	// Get returns the value stored for the key.
	// found=false if the key does not exist.
	Get(ctx context.Context, key string) (value string, found bool, err error)

	// Delete removes the key from the cache (idempotent).
	Delete(ctx context.Context, key string) error
}
