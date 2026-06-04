//go:build wireinject

package main

import (
	"zero-backend/internal/config"

	"zero-backend/modules/api/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() (*server.HTTPServer, error) {
	panic(wire.Build(
		config.New,
		providers.RedisProviderSet,
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
