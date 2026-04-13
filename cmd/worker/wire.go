//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/storage/redis"
	"zero-backend/modules/worker/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() *server.WorkerServer {
	panic(wire.Build(
		config.New,
		redis.New,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.WorkerProviderSet,
	))
}
