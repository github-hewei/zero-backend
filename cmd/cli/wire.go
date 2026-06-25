//go:build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/modules/cli/command"
	"zero-backend/providers"

	"github.com/241x/zero-kit/locker"

	"github.com/google/wire"
)

func wireApp() (*command.RootCommand, error) {
	panic(wire.Build(
		config.New,
		providers.RedisProviderSet,
		locker.NewRedisLocker,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.RepositoryProviderSet,
		providers.MySQLProviderSet,
		providers.CliCommandProviderSet,
	))
}
