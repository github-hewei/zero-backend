package providers

import (
	"zero-backend/internal/config"
)

// NewAdminAuthConfig 提取管理端认证配置
func NewAdminAuthConfig(cfg *config.Config) config.AdminAuthConfig {
	return cfg.Admin.Auth
}

// NewApiAuthConfig 提取API端认证配置
func NewApiAuthConfig(cfg *config.Config) config.ApiAuthConfig {
	return cfg.Api.Auth
}

// NewAdminCorsConfig 提取管理端跨域配置
func NewAdminCorsConfig(cfg *config.Config) config.AdminCorsConfig {
	return config.AdminCorsConfig(cfg.Admin.Cors)
}

// NewApiCorsConfig 提取API端跨域配置
func NewApiCorsConfig(cfg *config.Config) config.ApiCorsConfig {
	return config.ApiCorsConfig(cfg.Api.Cors)
}

// NewAdminServerConfig 提取管理端服务配置
func NewAdminServerConfig(cfg *config.Config) config.ServerConfig {
	return cfg.Admin.Server
}

// NewApiServerConfig 提取API端服务配置
func NewApiServerConfig(cfg *config.Config) config.ServerConfig {
	return cfg.Api.Server
}
