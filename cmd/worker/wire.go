//go:build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/modules/worker/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() *server.WorkerServer {
	panic(wire.Build(
		config.New,
		providers.RedisProviderSet,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.WorkerProviderSet,
	))
}
