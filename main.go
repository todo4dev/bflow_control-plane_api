package main

import (
	"context"
	"strconv"

	_ "src/application"
	_ "src/domain"
	_ "src/infrastructure"

	"src/application/adapter/logger"
	"src/application/usecase/system/query"
	"src/core/cqrs"
	"src/core/di"
	"src/core/env"
	"src/presentation/api"
)

func main() {
	env.Load("./.env", "../.env")

	cqrs.MustExecuteQuery[query.HealthQueryResult](context.Background(), &query.HealthQuery{})

	logger := di.Resolve[logger.ILoggerAdapter]()
	server := di.Resolve[*api.Server]()

	logger.Info("Server started on :" + strconv.Itoa(server.Config.Port))
	if err := server.ListenAndServe(context.Background()); err != nil {
		panic(err)
	}
}
