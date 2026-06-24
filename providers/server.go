package providers

import (
	"zero-backend/internal/config"
	"zero-backend/modules/admin/router"
	router2 "zero-backend/modules/api/router"

	config2 "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/server"
	"github.com/google/wire"
)

func ProvideServerOptions() []server.Option {
	return nil
}

// NewAdminServerConfig 提取管理端服务配置
func NewAdminServerConfig(cfg *config.Config) config2.ServerConfig {
	return cfg.Admin.Server
}

// NewApiServerConfig 提取API端服务配置
func NewApiServerConfig(cfg *config.Config) config2.ServerConfig {
	return cfg.Api.Server
}

// NewAdminCorsConfig 提取管理端跨域配置
func NewAdminCorsConfig(cfg *config.Config) config2.CorsConfig {
	return cfg.Admin.Cors
}

// NewApiCorsConfig 提取API端跨域配置
func NewApiCorsConfig(cfg *config.Config) config2.CorsConfig {
	return cfg.Api.Cors
}

// AdminServerProviderSet 提供管理端服务依赖集合
var AdminServerProviderSet = wire.NewSet(
	NewAdminServerConfig,
	NewAdminCorsConfig,
	ProvideServerOptions,
	server.New,
	router.NewGin,
)

// ApiServerProviderSet 提供API服务依赖集合
var ApiServerProviderSet = wire.NewSet(
	NewApiServerConfig,
	NewApiCorsConfig,
	ProvideServerOptions,
	server.New,
	router2.NewGin,
)
