package user

import "zero-backend/internal/config"

// Config API 端模块配置。
type Config struct {
	HmacSecret      string
	AccessTokenTtl  int
	RefreshTokenTtl int
}

// LoadConfig 从全局配置加载模块配置。
func LoadConfig() Config {
	var c Config
	config.UnmarshalKey("api.auth", &c)
	return c
}
