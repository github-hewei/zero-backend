package providers

import (
	"zero-backend/internal/config"
	"zero-backend/internal/constants"
	"zero-backend/modules/captcha"
	"zero-backend/modules/setting"
	service2 "zero-backend/modules/api/service"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewSettingService(db *gorm.DB) *setting.Service {
	return setting.NewService(setting.NewRepository(db), setting.NewDefaultRepository(db))
}

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

// NewCaptchaService 创建验证码服务
func NewCaptchaService(rdb *redis.Client, cfg config.CaptchaConfig) *captcha.Service {
	return captcha.NewService(rdb, captcha.Config{Enabled: cfg.Enabled, TTL: cfg.TTL}, constants.RedisCaptchaKey)
}

// AdminServiceProviderSet 提供管理端服务层依赖集合
var AdminServiceProviderSet = wire.NewSet(
	NewAdminAuthConfig,
	NewCaptchaConfig,
	NewCaptchaService,
	NewSettingService,
)

// ApiServiceProviderSet 提供API端服务层依赖集合
var ApiServiceProviderSet = wire.NewSet(NewApiAuthConfig, NewSettingService, service2.NewAuthService)
