package user

import (
	"errors"
	"zero-backend/internal/config"
)

// Config API 端模块配置。
type Config struct {
	HmacSecret      string
	AccessTokenTtl  int
	RefreshTokenTtl int
}

// Validate 校验配置有效性。
func (c Config) Validate() error {
	if c.HmacSecret == "" {
		return errors.New("user: api.auth.hmac_secret is required")
	}
	if c.AccessTokenTtl <= 0 {
		return errors.New("user: api.auth.access_token_ttl must be positive")
	}
	if c.RefreshTokenTtl <= 0 {
		return errors.New("user: api.auth.refresh_token_ttl must be positive")
	}
	return nil
}

// LoadConfig 从全局配置加载模块配置，校验失败返回 error。
func LoadConfig() (Config, error) {
	type adapter struct {
		HmacSecret      string `mapstructure:"hmac_secret"`
		AccessTokenTtl  int    `mapstructure:"access_token_ttl"`
		RefreshTokenTtl int    `mapstructure:"refresh_token_ttl"`
	}
	var a adapter
	config.UnmarshalKey("api.auth", &a)
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

func MustLoadConfig() Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
