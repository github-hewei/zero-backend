//go:build wireinject

package main

import (
	"zero-backend/config"
	"zero-backend/cmd/worker/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() (*server.WorkerServer, error) {
	panic(wire.Build(
		config.New,
		providers.RedisProviderSet,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.WorkerProviderSet,
	))
}
