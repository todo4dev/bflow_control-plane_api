package crypto

import (
	"src/core/validator"
)

type CryptoConfig struct {
	Key string
}

var _ validator.IValidable = (*CryptoConfig)(nil)

func (c *CryptoConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.Key).Required().Length(32),
	).Validate()
}
