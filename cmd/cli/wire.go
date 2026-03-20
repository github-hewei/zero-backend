//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/modules/cli/command"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireRootCommand() *command.RootCommand {
	panic(wire.Build(
		config.New,
		providers.LoggerProviderSet,
		providers.MongoDBProviderSet,
		providers.ServiceProviderSet,
		providers.RepositoryProviderSet,
		providers.MySQLProviderSet,
		providers.CliCommandProviderSet,
	))
}
