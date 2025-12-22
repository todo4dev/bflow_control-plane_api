package stream

import (
	"context"
	"time"
)

type Payload struct {
	Message    []byte         `json:"message,omitempty"`
	Key        *string        `json:"key,omitempty"`
	TTL        *time.Duration `json:"ttl,omitempty"`
	MaxRetries *int           `json:"max_retries,omitempty"`
}

type IStreamAdapter interface {
	Ping(ctx context.Context) error
	Publish(ctx context.Context, topic string, payload Payload) error
	Subscribe(ctx context.Context, topic string, handler func(payload Payload) error) error
}
