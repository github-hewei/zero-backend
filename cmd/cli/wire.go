//go:build wireinject

package main

import (
	"zero-backend/cmd/cli/command"
	"zero-backend/internal/config"
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
