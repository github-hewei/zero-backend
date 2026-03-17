//go:build wireinject
// +build wireinject

package main

import (
	"zero-backend/internal/config"
	"zero-backend/modules/cli"
	"zero-backend/providers"

	"github.com/google/wire"
)

func wireCLIContext() *cli.Context {
	panic(wire.Build(
		config.New,
		providers.LoggerProviderSet,
		providers.MySQLProviderSet,
		providers.MongoDBProviderSet,
		cli.NewContext,
	))
}
