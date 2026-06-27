package platform_user

import (
	"errors"
	"zero-backend/internal/config"
)

// Config 平台模块配置。
type Config struct {
	HmacSecret      string
	AccessTokenTtl  int
	RefreshTokenTtl int
}

// Validate 校验配置有效性。
func (c Config) Validate() error {
	if c.HmacSecret == "" {
		return errors.New("platform_user: platform.auth.hmac_secret is required")
	}
	if c.AccessTokenTtl <= 0 {
		return errors.New("platform_user: platform.auth.access_token_ttl must be positive")
	}
	if c.RefreshTokenTtl <= 0 {
		return errors.New("platform_user: platform.auth.refresh_token_ttl must be positive")
	}
	return nil
}

// LoadConfig 从全局配置加载模块配置，校验失败返回 error。
func LoadConfig() (Config, error) {
	c := Config{
		HmacSecret:      config.GetString("platform.auth.hmac_secret"),
		AccessTokenTtl:  config.GetInt("platform.auth.access_token_ttl"),
		RefreshTokenTtl: config.GetInt("platform.auth.refresh_token_ttl"),
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}
