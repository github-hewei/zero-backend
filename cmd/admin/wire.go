//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/storage/redis"

	"zero-backend/modules/admin/server"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() *server.HTTPServer {
	panic(wire.Build(
		config.New,
		redis.New,
		providers.AdminControllersProviderSet,
		providers.MiddlewaresProviderSet,
		providers.AdminMiddlewaresProviderSet,
		providers.RepositoryProviderSet,
		providers.ServiceProviderSet,
		providers.AdminServiceProviderSet,
		providers.AdminServerProviderSet,
		providers.RequestProviderSet,
		providers.LoggerProviderSet,
		providers.MySQLProviderSet,
		providers.MongoDBProviderSet,
	))
}
