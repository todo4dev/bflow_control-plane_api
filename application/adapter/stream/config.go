package stream

import "src/core/validator"

type StreamConfig struct {
	KafkaURI string
}

var _ validator.IValidable = (*StreamConfig)(nil)

func (c *StreamConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.KafkaURI).Required().Default("localhost:9092"),
	).Validate()
}
