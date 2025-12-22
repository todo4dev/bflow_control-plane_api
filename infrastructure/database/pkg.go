package database

import (
	"context"

	adapter "src/application/adapter/database"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/database/pgx"
)

func init() {
	di.RegisterAs[adapter.IDatabaseAdapter](func() adapter.IDatabaseAdapter {
		config := &adapter.DatabaseConfig{
			PostgresURI:  env.Get("DATABASE_POSTGRES_URI", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
			DefaultLimit: env.Get[int64]("DATABASE_DEFAULT_LIMIT", 100),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		pgx := impl.NewPgxDatabaseAdapter(config)
		if err := pgx.Ping(context.Background()); err != nil {
			panic(err)
		}
		return pgx
	})
}
