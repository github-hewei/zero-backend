package rbac

import (
	"errors"
	"zero-backend/internal/config"
)

// Config 模块配置。
type Config struct {
	HmacSecret      string
	AccessTokenTtl  int
	RefreshTokenTtl int
	SuperUserId     int
}

// Validate 校验配置有效性。
func (c Config) Validate() error {
	if c.HmacSecret == "" {
		return errors.New("rbac: admin.auth.hmac_secret is required")
	}
	if c.AccessTokenTtl <= 0 {
		return errors.New("rbac: admin.auth.access_token_ttl must be positive")
	}
	if c.RefreshTokenTtl <= 0 {
		return errors.New("rbac: admin.auth.refresh_token_ttl must be positive")
	}
	return nil
}

// LoadConfig 从全局配置加载模块配置，校验失败返回 error。
func LoadConfig() (Config, error) {
	var c Config
	config.UnmarshalKey("admin.auth", &c)
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}
