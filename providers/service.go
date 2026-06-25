package providers

import (
	"zero-backend/internal/config"
	"zero-backend/internal/service"
	service2 "zero-backend/modules/admin/service"
	service3 "zero-backend/modules/api/service"

	"github.com/google/wire"
)

// ServiceProviderSet 提供服务层依赖集合
var ServiceProviderSet = wire.NewSet(
	service.NewRbacMenuService,
	service.NewRbacApiService,
	service.NewRbacRoleService,
	service.NewRbacUserService,
	service.NewRbacStoreService,
	service.NewUserService,
	service.NewSettingService,
	service.NewSettingDefaultService,
	service.NewCaptchaService,
)

// NewAdminAuthConfig 提取管理端认证配置
func NewAdminAuthConfig(cfg *config.Config) config.AdminAuthConfig {
	return cfg.Admin.Auth
}

// NewCaptchaConfig 提取验证码配置
func NewCaptchaConfig(cfg *config.Config) config.CaptchaConfig {
	return cfg.Admin.Captcha
}

// NewApiAuthConfig 提取 API端认证配置
func NewApiAuthConfig(cfg *config.Config) config.ApiAuthConfig {
	return cfg.Api.Auth
}

// AdminServiceProviderSet 提供管理端服务层依赖集合
var AdminServiceProviderSet = wire.NewSet(
	NewAdminAuthConfig,
	NewCaptchaConfig,
	service2.NewAuthService,
)

// ApiServiceProviderSet 提供API端服务层依赖集合
var ApiServiceProviderSet = wire.NewSet(NewApiAuthConfig, service3.NewAuthService)
