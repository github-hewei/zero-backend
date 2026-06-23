package providers

import (
	"zero-backend/internal/config"
	"zero-backend/modules/admin/server"
	server2 "zero-backend/modules/api/server"

	baseconfig "github.com/241x/zero-web/config"
	"github.com/google/wire"
)

// NewAdminServerConfig 提取管理端服务配置
func NewAdminServerConfig(cfg *config.Config) baseconfig.ServerConfig {
	return cfg.Admin.Server
}

// NewApiServerConfig 提取API端服务配置
func NewApiServerConfig(cfg *config.Config) baseconfig.ServerConfig {
	return cfg.Api.Server
}

// NewAdminCorsConfig 提取管理端跨域配置
func NewAdminCorsConfig(cfg *config.Config) baseconfig.CorsConfig {
	return cfg.Admin.Cors
}

// NewApiCorsConfig 提取API端跨域配置
func NewApiCorsConfig(cfg *config.Config) baseconfig.CorsConfig {
	return cfg.Api.Cors
}

// AdminServerProviderSet 提供管理端服务依赖集合
var AdminServerProviderSet = wire.NewSet(
	NewAdminServerConfig,
	NewAdminCorsConfig,
	server.NewHTTPServer,
	server.NewGin,
)

// ApiServerProviderSet 提供API服务依赖集合
var ApiServerProviderSet = wire.NewSet(
	NewApiServerConfig,
	NewApiCorsConfig,
	server2.NewHTTPServer,
	server2.NewGin,
)
