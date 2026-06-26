package app

import (
	"zero-backend/internal/config"

	baseconfig "github.com/241x/zero-web/config"
	webserver "github.com/241x/zero-web/server"
)

func ProvideServerOptions() []webserver.Option { return nil }

func NewAdminServerConfig(cfg *config.Config) baseconfig.ServerConfig { return cfg.Admin.Server }
func NewApiServerConfig(cfg *config.Config) baseconfig.ServerConfig   { return cfg.Api.Server }
func NewAdminCorsConfig(cfg *config.Config) baseconfig.CorsConfig     { return cfg.Admin.Cors }
func NewApiCorsConfig(cfg *config.Config) baseconfig.CorsConfig       { return cfg.Api.Cors }
