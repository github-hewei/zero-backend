//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/storage/redis"

	"zero-backend/modules/api/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() *server.HTTPServer {
	panic(wire.Build(
		config.New,
		redis.New,
		providers.ApiControllersProviderSet,
		providers.MiddlewaresProviderSet,
		providers.ApiMiddlewaresProviderSet,
		providers.RepositoryProviderSet,
		providers.ServiceProviderSet,
		providers.ApiServiceProviderSet,
		providers.ApiServerProviderSet,
		providers.RequestProviderSet,
		providers.LoggerProviderSet,
		providers.MySQLProviderSet,
		providers.MongoDBProviderSet,
	))
}
