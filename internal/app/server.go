package app

import (
	"zero-backend/internal/config"

	baseconfig "github.com/241x/zero-web/config"
	webserver "github.com/241x/zero-web/server"
)

func ProvideServerOptions() []webserver.Option { return nil }

func LoadAdminServerConfig() baseconfig.ServerConfig {
	var c baseconfig.ServerConfig
	config.UnmarshalKey("admin.server", &c)
	return c
}

func LoadApiServerConfig() baseconfig.ServerConfig {
	var c baseconfig.ServerConfig
	config.UnmarshalKey("api.server", &c)
	return c
}

func LoadAdminCorsConfig() baseconfig.CorsConfig {
	var c baseconfig.CorsConfig
	config.UnmarshalKey("admin.cors", &c)
	return c
}

func LoadApiCorsConfig() baseconfig.CorsConfig {
	var c baseconfig.CorsConfig
	config.UnmarshalKey("api.cors", &c)
	return c
}
