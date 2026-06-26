package app

import (
	"zero-backend/internal/config"
	"zero-backend/internal/modules/captcha"
	"zero-backend/internal/modules/setting"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewSettingService(db *gorm.DB) *setting.Service {
	return setting.NewService(setting.NewRepository(db), setting.NewDefaultRepository(db))
}

func NewAdminAuthConfig(cfg *config.Config) config.AdminAuthConfig { return cfg.Admin.Auth }
func NewCaptchaConfig(cfg *config.Config) config.CaptchaConfig     { return cfg.Admin.Captcha }
func NewApiAuthConfig(cfg *config.Config) config.ApiAuthConfig     { return cfg.Api.Auth }

func NewCaptchaService(rdb *goredis.Client, cfg config.CaptchaConfig) *captcha.Service {
	return captcha.NewService(rdb, captcha.Config{Enabled: cfg.Enabled, TTL: cfg.TTL}, "ZAG:CAPTCHA")
}
