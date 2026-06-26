package app

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-web/middleware"
	"github.com/241x/zero-web/server"
)

func ProvideServerOptions() []server.Option { return nil }

func LoadAdminServerConfig() server.Config {
	return server.Config{
		Host: config.GetString("admin.server.host"),
		Port: config.GetInt("admin.server.port"),
	}
}

func LoadApiServerConfig() server.Config {
	return server.Config{
		Host: config.GetString("api.server.host"),
		Port: config.GetInt("api.server.port"),
	}
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
