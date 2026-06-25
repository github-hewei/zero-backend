//go:build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/providers"

	webserver "github.com/241x/zero-web/server"
	"github.com/google/wire"
)

func wireApp() (*webserver.Server, error) {
	panic(wire.Build(
		config.New,
		providers.RedisProviderSet,
		providers.ApiControllersProviderSet,
		providers.ApiMiddlewaresProviderSet,
		providers.RepositoryProviderSet,
		providers.ApiServiceProviderSet,
		providers.ApiServerProviderSet,
		providers.BindProviderSet,
		providers.LoggerProviderSet,
		providers.MySQLProviderSet,
		providers.MongoDBProviderSet,
	))
}
