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

func LoadCaptchaConfig() captcha.Config {
	var c struct {
		Enabled bool
		TTL     int
	}
	config.UnmarshalKey("admin.captcha", &c)
	return captcha.Config{Enabled: c.Enabled, TTL: c.TTL}
}

func NewCaptchaService(rdb *goredis.Client, cfg captcha.Config) (*captcha.Service, error) {
	return captcha.NewService(rdb, cfg, "ZAG:CAPTCHA")
}
