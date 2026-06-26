package app

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/server"
)

func ProvideServerOptions() []server.Option { return nil }

func LoadAdminServerConfig() server.Config {
	var c server.Config
	config.UnmarshalKey("admin.server", &c)
	return c
}

func LoadApiServerConfig() server.Config {
	var c server.Config
	config.UnmarshalKey("api.server", &c)
	return c
}

func LoadAdminCorsConfig() middleware.CorsConfig {
	var c middleware.CorsConfig
	config.UnmarshalKey("admin.cors", &c)
	return c
}

func LoadApiCorsConfig() middleware.CorsConfig {
	var c middleware.CorsConfig
	config.UnmarshalKey("api.cors", &c)
	return c
}
