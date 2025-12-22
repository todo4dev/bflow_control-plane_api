package stream

import (
	"context"

	adapter "src/application/adapter/stream"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/stream/kafka"
)

func init() {
	di.RegisterAs[adapter.IStreamAdapter](func() adapter.IStreamAdapter {
		config := &adapter.StreamConfig{
			KafkaURI: env.Get("STREAM_KAFKA_URI", "localhost:9092"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		impl := impl.NewKafkaStreamAdapter(config)
		if err := impl.Ping(context.Background()); err != nil {
			panic(err)
		}
		return impl
	})
}
