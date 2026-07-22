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

// LoadAdminConfig 从全局配置加载模块配置，校验失败返回 error。
func LoadAdminConfig() (Config, error) {
	type adapter struct {
		HmacSecret      string `mapstructure:"hmac_secret"`
		AccessTokenTtl  int    `mapstructure:"access_token_ttl"`
		RefreshTokenTtl int    `mapstructure:"refresh_token_ttl"`
	}
	var a adapter
	config.UnmarshalKey("admin.auth", &a)
	c := Config{
		HmacSecret:      a.HmacSecret,
		AccessTokenTtl:  a.AccessTokenTtl,
		RefreshTokenTtl: a.RefreshTokenTtl,
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

// MustLoadAdminConfig 从全局配置加载模块配置，校验失败 panic。
func MustLoadAdminConfig() Config {
	cfg, err := LoadAdminConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

// LoadPlatformConfig 从全局配置加载模块配置，校验失败返回 error。
func LoadPlatformConfig() (Config, error) {
	type adapter struct {
		HmacSecret      string `mapstructure:"hmac_secret"`
		AccessTokenTtl  int    `mapstructure:"access_token_ttl"`
		RefreshTokenTtl int    `mapstructure:"refresh_token_ttl"`
	}
	var a adapter
	config.UnmarshalKey("platform.auth", &a)
	c := Config{
		HmacSecret:      a.HmacSecret,
		AccessTokenTtl:  a.AccessTokenTtl,
		RefreshTokenTtl: a.RefreshTokenTtl,
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

// MustLoadPlatformConfig 从全局配置加载模块配置，校验失败 panic。
func MustLoadPlatformConfig() Config {
	cfg, err := LoadPlatformConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
