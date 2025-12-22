package database

import (
	"src/core/validator"
)

type DatabaseConfig struct {
	PostgresURI  string
	DefaultLimit int64
}

var _ validator.IValidable = (*DatabaseConfig)(nil)

func (c *DatabaseConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.PostgresURI).Required().URI().Default("postgres://postgres:postgres@localhost:5432/postgres"),
		validator.Number(&c.DefaultLimit).Required().Default(50),
	).Validate()
}
