//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/internal/storage/redis"
	"zero-backend/modules/cli/command"
	"zero-backend/pkg/locker"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireApp() *command.RootCommand {
	panic(wire.Build(
		config.New,
		redis.New,
		locker.NewRedisLocker,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.ServiceProviderSet,
		providers.RepositoryProviderSet,
		providers.MySQLProviderSet,
		providers.CliCommandProviderSet,
	))
}
